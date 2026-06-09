package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type audioFilePlanModifier struct{}

func (audioFilePlanModifier) Description(_ context.Context) string {
	return "Recomputes uploaded audio metadata when the source file changes."
}

func (audioFilePlanModifier) MarkdownDescription(_ context.Context) string {
	return "Recomputes uploaded audio metadata when the source file changes."
}

func (m audioFilePlanModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	planAttrs := req.PlanValue.Attributes()
	sourceValue, ok := planAttrs["source"].(types.String)
	if !ok || sourceValue.IsNull() || sourceValue.IsUnknown() {
		return
	}

	newHash, err := hashFile(sourceValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to hash audio source file", err.Error())
		return
	}

	if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
		stateAttrs := req.StateValue.Attributes()
		if stateHash, ok := stateAttrs["source_hash"].(types.String); ok && stateHash.ValueString() == newHash {
			resp.PlanValue = req.StateValue
			return
		}
	}

	planValue, diags := types.ObjectValue(
		map[string]attr.Type{
			"source":          types.StringType,
			"object_key":      types.StringType,
			"file_name":       types.StringType,
			"content_type":    types.StringType,
			"file_size_bytes": types.Int64Type,
			"source_hash":     types.StringType,
		},
		map[string]attr.Value{
			"source":          sourceValue,
			"object_key":      types.StringUnknown(),
			"file_name":       types.StringUnknown(),
			"content_type":    types.StringUnknown(),
			"file_size_bytes": types.Int64Unknown(),
			"source_hash":     types.StringUnknown(),
		},
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.PlanValue = planValue
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("read %s: %w", path, err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func uploadChangedTwilioAudio(ctx context.Context, client *AllQuietAPIClient, integrationId string, settings *TwilioSettingsModel) error {
	if settings == nil || settings.CallFlowConfig == nil {
		return nil
	}

	config := settings.CallFlowConfig
	if err := uploadPromptAudioIfNeeded(ctx, client, integrationId, config.WelcomeMessage, callRoutingMediaPurposeWelcome); err != nil {
		return err
	}
	if config.Route != nil {
		if err := uploadRoutePromptsAudio(ctx, client, integrationId, config.Route); err != nil {
			return err
		}
	}
	if config.Menu != nil {
		for i := range config.Menu.Options {
			if err := uploadRoutePromptsAudio(ctx, client, integrationId, &config.Menu.Options[i].Route); err != nil {
				return err
			}
		}
	}
	return nil
}

func uploadRoutePromptsAudio(ctx context.Context, client *AllQuietAPIClient, integrationId string, route *RouteConfigModel) error {
	if route == nil {
		return nil
	}
	if err := uploadPromptAudioIfNeeded(ctx, client, integrationId, route.VoicemailMessage, callRoutingMediaPurposeVoicemail); err != nil {
		return err
	}
	return uploadPromptAudioIfNeeded(ctx, client, integrationId, route.ClosingMessage, callRoutingMediaPurposeClosing)
}

func uploadPromptAudioIfNeeded(ctx context.Context, client *AllQuietAPIClient, integrationId string, prompt *PromptModel, purpose string) error {
	if prompt == nil || prompt.AudioFile == nil {
		return nil
	}
	if prompt.Type.ValueString() != "Audio" {
		return nil
	}

	audio := prompt.AudioFile
	if audio.Source.IsNull() || audio.Source.IsUnknown() {
		return nil
	}
	if !audioNeedsUpload(audio) {
		return nil
	}

	uploaded, err := client.UploadCallRoutingMedia(ctx, integrationId, purpose, audio.Source.ValueString())
	if err != nil {
		return err
	}

	audio.ObjectKey = types.StringValue(uploaded.ObjectKey)
	audio.FileName = types.StringPointerValue(uploaded.FileName)
	audio.ContentType = types.StringPointerValue(uploaded.ContentType)
	if uploaded.FileSizeBytes != nil {
		audio.FileSizeBytes = types.Int64Value(*uploaded.FileSizeBytes)
	}
	if hash, hashErr := hashFile(audio.Source.ValueString()); hashErr == nil {
		audio.SourceHash = types.StringValue(hash)
	}

	return nil
}

func twilioSettingsNeedsAudioUpload(settings *TwilioSettingsModel) bool {
	if settings == nil || settings.CallFlowConfig == nil {
		return false
	}
	return callFlowNeedsAudioUpload(settings.CallFlowConfig)
}

func callFlowNeedsAudioUpload(config *CallFlowConfigModel) bool {
	if promptNeedsAudioUpload(config.WelcomeMessage) {
		return true
	}
	if config.Route != nil && routeNeedsAudioUpload(config.Route) {
		return true
	}
	if config.Menu != nil {
		for i := range config.Menu.Options {
			if routeNeedsAudioUpload(&config.Menu.Options[i].Route) {
				return true
			}
		}
	}
	return false
}

func routeNeedsAudioUpload(route *RouteConfigModel) bool {
	return promptNeedsAudioUpload(route.VoicemailMessage) || promptNeedsAudioUpload(route.ClosingMessage)
}

func promptNeedsAudioUpload(prompt *PromptModel) bool {
	if prompt == nil || prompt.AudioFile == nil {
		return false
	}
	if prompt.Type.ValueString() != "Audio" {
		return false
	}
	return audioNeedsUpload(prompt.AudioFile)
}

func audioNeedsUpload(audio *AudioFileModel) bool {
	if audio.Source.IsNull() || audio.Source.IsUnknown() {
		return false
	}
	if audio.ObjectKey.IsUnknown() || audio.ObjectKey.ValueString() == "" {
		return true
	}

	newHash, err := hashFile(audio.Source.ValueString())
	if err != nil {
		return true
	}
	if audio.SourceHash.IsNull() || audio.SourceHash.IsUnknown() {
		return true
	}

	return audio.SourceHash.ValueString() != newHash
}
