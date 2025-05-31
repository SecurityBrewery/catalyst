package testing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app2"
	"github.com/SecurityBrewery/catalyst/reaction"
)

func App(t *testing.T) (*app2.App2, *Counter) {
	t.Helper()

	temp := t.TempDir()

	baseApp, err := app2.App(t.Context(), temp, true)
	require.NoError(t, err)

	err = baseApp.SetupRoutes()
	require.NoError(t, err)

	reaction.BindHooks(baseApp, true)

	defaultTestData(t, baseApp)

	counter := countEvents(baseApp)

	return baseApp, counter
}

func countEvents(app *app2.App2) *Counter {
	c := NewCounter()

	// app.Service.OnBeforeApiError().Add(count[*core.ApiErrorEvent](c, "OnBeforeApiError"))
	// app.Service.OnBeforeApiError().Add(count[*core.ApiErrorEvent](c, "OnBeforeApiError"))
	// app.Service.OnAfterApiError().Add(count[*core.ApiErrorEvent](c, "OnAfterApiError"))
	// app.Service.OnModelBeforeCreate().Add(count[*core.ModelEvent](c, "OnModelBeforeCreate"))
	// app.Service.OnModelAfterCreate().Add(count[*core.ModelEvent](c, "OnModelAfterCreate"))
	// app.Service.OnModelBeforeUpdate().Add(count[*core.ModelEvent](c, "OnModelBeforeUpdate"))
	// app.Service.OnModelAfterUpdate().Add(count[*core.ModelEvent](c, "OnModelAfterUpdate"))
	// app.Service.OnModelBeforeDelete().Add(count[*core.ModelEvent](c, "OnModelBeforeDelete"))
	// app.Service.OnModelAfterDelete().Add(count[*core.ModelEvent](c, "OnModelAfterDelete"))
	app.Service.OnRecordsListRequest.Subscribe(count(c, "OnRecordsListRequest"))
	app.Service.OnRecordViewRequest.Subscribe(count(c, "OnRecordViewRequest"))
	app.Service.OnRecordBeforeCreateRequest.Subscribe(count(c, "OnRecordBeforeCreateRequest"))
	app.Service.OnRecordAfterCreateRequest.Subscribe(count(c, "OnRecordAfterCreateRequest"))
	app.Service.OnRecordBeforeUpdateRequest.Subscribe(count(c, "OnRecordBeforeUpdateRequest"))
	app.Service.OnRecordAfterUpdateRequest.Subscribe(count(c, "OnRecordAfterUpdateRequest"))
	app.Service.OnRecordBeforeDeleteRequest.Subscribe(count(c, "OnRecordBeforeDeleteRequest"))
	app.Service.OnRecordAfterDeleteRequest.Subscribe(count(c, "OnRecordAfterDeleteRequest"))
	// app.Service.OnRecordAuthRequest().Add(count[*core.RecordAuthEvent](c, "OnRecordAuthRequest"))
	// app.Service.OnRecordBeforeAuthWithPasswordRequest().Add(count[*core.RecordAuthWithPasswordEvent](c, "OnRecordBeforeAuthWithPasswordRequest"))
	// app.Service.OnRecordAfterAuthWithPasswordRequest().Add(count[*core.RecordAuthWithPasswordEvent](c, "OnRecordAfterAuthWithPasswordRequest"))
	// app.Service.OnRecordBeforeAuthWithOAuth2Request().Add(count[*core.RecordAuthWithOAuth2Event](c, "OnRecordBeforeAuthWithOAuth2Request"))
	// app.Service.OnRecordAfterAuthWithOAuth2Request().Add(count[*core.RecordAuthWithOAuth2Event](c, "OnRecordAfterAuthWithOAuth2Request"))
	// app.Service.OnRecordBeforeAuthRefreshRequest().Add(count[*core.RecordAuthRefreshEvent](c, "OnRecordBeforeAuthRefreshRequest"))
	// app.Service.OnRecordAfterAuthRefreshRequest().Add(count[*core.RecordAuthRefreshEvent](c, "OnRecordAfterAuthRefreshRequest"))
	// app.Service.OnRecordBeforeRequestPasswordResetRequest().Add(count[*core.RecordRequestPasswordResetEvent](c, "OnRecordBeforeRequestPasswordResetRequest"))
	// app.Service.OnRecordAfterRequestPasswordResetRequest().Add(count[*core.RecordRequestPasswordResetEvent](c, "OnRecordAfterRequestPasswordResetRequest"))
	// app.Service.OnRecordBeforeConfirmPasswordResetRequest().Add(count[*core.RecordConfirmPasswordResetEvent](c, "OnRecordBeforeConfirmPasswordResetRequest"))
	// app.Service.OnRecordAfterConfirmPasswordResetRequest().Add(count[*core.RecordConfirmPasswordResetEvent](c, "OnRecordAfterConfirmPasswordResetRequest"))
	// app.Service.OnRecordBeforeRequestVerificationRequest().Add(count[*core.RecordRequestVerificationEvent](c, "OnRecordBeforeRequestVerificationRequest"))
	// app.Service.OnRecordAfterRequestVerificationRequest().Add(count[*core.RecordRequestVerificationEvent](c, "OnRecordAfterRequestVerificationRequest"))
	// app.Service.OnRecordBeforeConfirmVerificationRequest().Add(count[*core.RecordConfirmVerificationEvent](c, "OnRecordBeforeConfirmVerificationRequest"))
	// app.Service.OnRecordAfterConfirmVerificationRequest().Add(count[*core.RecordConfirmVerificationEvent](c, "OnRecordAfterConfirmVerificationRequest"))
	// app.Service.OnRecordBeforeRequestEmailChangeRequest().Add(count[*core.RecordRequestEmailChangeEvent](c, "OnRecordBeforeRequestEmailChangeRequest"))
	// app.Service.OnRecordAfterRequestEmailChangeRequest().Add(count[*core.RecordRequestEmailChangeEvent](c, "OnRecordAfterRequestEmailChangeRequest"))
	// app.Service.OnRecordBeforeConfirmEmailChangeRequest().Add(count[*core.RecordConfirmEmailChangeEvent](c, "OnRecordBeforeConfirmEmailChangeRequest"))
	// app.Service.OnRecordAfterConfirmEmailChangeRequest().Add(count[*core.RecordConfirmEmailChangeEvent](c, "OnRecordAfterConfirmEmailChangeRequest"))
	// app.Service.OnRecordListExternalAuthsRequest().Add(count[*core.RecordListExternalAuthsEvent](c, "OnRecordListExternalAuthsRequest"))
	// app.Service.OnRecordBeforeUnlinkExternalAuthRequest().Add(count[*core.RecordUnlinkExternalAuthEvent](c, "OnRecordBeforeUnlinkExternalAuthRequest"))
	// app.Service.OnRecordAfterUnlinkExternalAuthRequest().Add(count[*core.RecordUnlinkExternalAuthEvent](c, "OnRecordAfterUnlinkExternalAuthRequest"))
	// app.Service.OnMailerBeforeAdminResetPasswordSend().Add(count[*core.MailerAdminEvent](c, "OnMailerBeforeAdminResetPasswordSend"))
	// app.Service.OnMailerAfterAdminResetPasswordSend().Add(count[*core.MailerAdminEvent](c, "OnMailerAfterAdminResetPasswordSend"))
	// app.Service.OnMailerBeforeRecordResetPasswordSend().Add(count[*core.MailerRecordEvent](c, "OnMailerBeforeRecordResetPasswordSend"))
	// app.Service.OnMailerAfterRecordResetPasswordSend().Add(count[*core.MailerRecordEvent](c, "OnMailerAfterRecordResetPasswordSend"))
	// app.Service.OnMailerBeforeRecordVerificationSend().Add(count[*core.MailerRecordEvent](c, "OnMailerBeforeRecordVerificationSend"))
	// app.Service.OnMailerAfterRecordVerificationSend().Add(count[*core.MailerRecordEvent](c, "OnMailerAfterRecordVerificationSend"))
	// app.Service.OnMailerBeforeRecordChangeEmailSend().Add(count[*core.MailerRecordEvent](c, "OnMailerBeforeRecordChangeEmailSend"))
	// app.Service.OnMailerAfterRecordChangeEmailSend().Add(count[*core.MailerRecordEvent](c, "OnMailerAfterRecordChangeEmailSend"))
	// app.Service.OnRealtimeConnectRequest().Add(count[*core.RealtimeConnectEvent](c, "OnRealtimeConnectRequest"))
	// app.Service.OnRealtimeDisconnectRequest().Add(count[*core.RealtimeDisconnectEvent](c, "OnRealtimeDisconnectRequest"))
	// app.Service.OnRealtimeBeforeMessageSend().Add(count[*core.RealtimeMessageEvent](c, "OnRealtimeBeforeMessageSend"))
	// app.Service.OnRealtimeAfterMessageSend().Add(count[*core.RealtimeMessageEvent](c, "OnRealtimeAfterMessageSend"))
	// app.Service.OnRealtimeBeforeSubscribeRequest().Add(count[*core.RealtimeSubscribeEvent](c, "OnRealtimeBeforeSubscribeRequest"))
	// app.Service.OnRealtimeAfterSubscribeRequest().Add(count[*core.RealtimeSubscribeEvent](c, "OnRealtimeAfterSubscribeRequest"))
	// app.Service.OnSettingsListRequest().Add(count[*core.SettingsListEvent](c, "OnSettingsListRequest"))
	// app.Service.OnSettingsBeforeUpdateRequest().Add(count[*core.SettingsUpdateEvent](c, "OnSettingsBeforeUpdateRequest"))
	// app.Service.OnSettingsAfterUpdateRequest().Add(count[*core.SettingsUpdateEvent](c, "OnSettingsAfterUpdateRequest"))
	// app.Service.OnCollectionsListRequest().Add(count[*core.CollectionsListEvent](c, "OnCollectionsListRequest"))
	// app.Service.OnCollectionViewRequest().Add(count[*core.CollectionViewEvent](c, "OnCollectionViewRequest"))
	// app.Service.OnCollectionBeforeCreateRequest().Add(count[*core.CollectionCreateEvent](c, "OnCollectionBeforeCreateRequest"))
	// app.Service.OnCollectionAfterCreateRequest().Add(count[*core.CollectionCreateEvent](c, "OnCollectionAfterCreateRequest"))
	// app.Service.OnCollectionBeforeUpdateRequest().Add(count[*core.CollectionUpdateEvent](c, "OnCollectionBeforeUpdateRequest"))
	// app.Service.OnCollectionAfterUpdateRequest().Add(count[*core.CollectionUpdateEvent](c, "OnCollectionAfterUpdateRequest"))
	// app.Service.OnCollectionBeforeDeleteRequest().Add(count[*core.CollectionDeleteEvent](c, "OnCollectionBeforeDeleteRequest"))
	// app.Service.OnCollectionAfterDeleteRequest().Add(count[*core.CollectionDeleteEvent](c, "OnCollectionAfterDeleteRequest"))
	// app.Service.OnCollectionsBeforeImportRequest().Add(count[*core.CollectionsImportEvent](c, "OnCollectionsBeforeImportRequest"))
	// app.Service.OnCollectionsAfterImportRequest().Add(count[*core.CollectionsImportEvent](c, "OnCollectionsAfterImportRequest"))
	// app.Service.OnAdminsListRequest().Add(count[*core.AdminsListEvent](c, "OnAdminsListRequest"))
	// app.Service.OnAdminViewRequest().Add(count[*core.AdminViewEvent](c, "OnAdminViewRequest"))
	// app.Service.OnAdminBeforeCreateRequest().Add(count[*core.AdminCreateEvent](c, "OnAdminBeforeCreateRequest"))
	// app.Service.OnAdminAfterCreateRequest().Add(count[*core.AdminCreateEvent](c, "OnAdminAfterCreateRequest"))
	// app.Service.OnAdminBeforeUpdateRequest().Add(count[*core.AdminUpdateEvent](c, "OnAdminBeforeUpdateRequest"))
	// app.Service.OnAdminAfterUpdateRequest().Add(count[*core.AdminUpdateEvent](c, "OnAdminAfterUpdateRequest"))
	// app.Service.OnAdminBeforeDeleteRequest().Add(count[*core.AdminDeleteEvent](c, "OnAdminBeforeDeleteRequest"))
	// app.Service.OnAdminAfterDeleteRequest().Add(count[*core.AdminDeleteEvent](c, "OnAdminAfterDeleteRequest"))
	// app.Service.OnAdminAuthRequest().Add(count[*core.AdminAuthEvent](c, "OnAdminAuthRequest"))
	// app.Service.OnAdminBeforeAuthWithPasswordRequest().Add(count[*core.AdminAuthWithPasswordEvent](c, "OnAdminBeforeAuthWithPasswordRequest"))
	// app.Service.OnAdminAfterAuthWithPasswordRequest().Add(count[*core.AdminAuthWithPasswordEvent](c, "OnAdminAfterAuthWithPasswordRequest"))
	// app.Service.OnAdminBeforeAuthRefreshRequest().Add(count[*core.AdminAuthRefreshEvent](c, "OnAdminBeforeAuthRefreshRequest"))
	// app.Service.OnAdminAfterAuthRefreshRequest().Add(count[*core.AdminAuthRefreshEvent](c, "OnAdminAfterAuthRefreshRequest"))
	// app.Service.OnAdminBeforeRequestPasswordResetRequest().Add(count[*core.AdminRequestPasswordResetEvent](c, "OnAdminBeforeRequestPasswordResetRequest"))
	// app.Service.OnAdminAfterRequestPasswordResetRequest().Add(count[*core.AdminRequestPasswordResetEvent](c, "OnAdminAfterRequestPasswordResetRequest"))
	// app.Service.OnAdminBeforeConfirmPasswordResetRequest().Add(count[*core.AdminConfirmPasswordResetEvent](c, "OnAdminBeforeConfirmPasswordResetRequest"))
	// app.Service.OnAdminAfterConfirmPasswordResetRequest().Add(count[*core.AdminConfirmPasswordResetEvent](c, "OnAdminAfterConfirmPasswordResetRequest"))
	// app.Service.OnFileDownloadRequest().Add(count[*core.FileDownloadEvent](c, "OnFileDownloadRequest"))
	// app.Service.OnFileBeforeTokenRequest().Add(count[*core.FileTokenEvent](c, "OnFileBeforeTokenRequest"))
	// app.Service.OnFileAfterTokenRequest().Add(count[*core.FileTokenEvent](c, "OnFileAfterTokenRequest"))
	// app.Service.OnFileAfterTokenRequest().Add(count[*core.FileTokenEvent](c, "OnFileAfterTokenRequest"))

	return c
}

func count(c *Counter, name string) func(ctx context.Context, table string, record any) {
	return func(_ context.Context, _ string, _ any) {
		c.Increment(name)
	}
}
