package testing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/reaction"
)

func App(t *testing.T) (*app.App, func(), *Counter) {
	t.Helper()

	baseApp, cleanup, err := app.New(t.Context(), t.TempDir())
	require.NoError(t, err)

	err = baseApp.SetupRoutes()
	require.NoError(t, err)

	err = reaction.BindHooks(baseApp, true)
	require.NoError(t, err)

	database.DefaultTestData(t, baseApp.Queries)

	counter := countEvents(baseApp)

	return baseApp, cleanup, counter
}

func countEvents(app *app.App) *Counter {
	c := NewCounter()

	// app.Hooks.OnBeforeApiError().Add(count[*core.ApiErrorEvent](c, "OnBeforeApiError"))
	// app.Hooks.OnBeforeApiError().Add(count[*core.ApiErrorEvent](c, "OnBeforeApiError"))
	// app.Hooks.OnAfterApiError().Add(count[*core.ApiErrorEvent](c, "OnAfterApiError"))
	// app.Hooks.OnModelBeforeCreate().Add(count[*core.ModelEvent](c, "OnModelBeforeCreate"))
	// app.Hooks.OnModelAfterCreate().Add(count[*core.ModelEvent](c, "OnModelAfterCreate"))
	// app.Hooks.OnModelBeforeUpdate().Add(count[*core.ModelEvent](c, "OnModelBeforeUpdate"))
	// app.Hooks.OnModelAfterUpdate().Add(count[*core.ModelEvent](c, "OnModelAfterUpdate"))
	// app.Hooks.OnModelBeforeDelete().Add(count[*core.ModelEvent](c, "OnModelBeforeDelete"))
	// app.Hooks.OnModelAfterDelete().Add(count[*core.ModelEvent](c, "OnModelAfterDelete"))
	app.Hooks.OnRecordsListRequest.Subscribe(count(c, "OnRecordsListRequest"))
	app.Hooks.OnRecordViewRequest.Subscribe(count(c, "OnRecordViewRequest"))
	app.Hooks.OnRecordBeforeCreateRequest.Subscribe(count(c, "OnRecordBeforeCreateRequest"))
	app.Hooks.OnRecordAfterCreateRequest.Subscribe(count(c, "OnRecordAfterCreateRequest"))
	app.Hooks.OnRecordBeforeUpdateRequest.Subscribe(count(c, "OnRecordBeforeUpdateRequest"))
	app.Hooks.OnRecordAfterUpdateRequest.Subscribe(count(c, "OnRecordAfterUpdateRequest"))
	app.Hooks.OnRecordBeforeDeleteRequest.Subscribe(count(c, "OnRecordBeforeDeleteRequest"))
	app.Hooks.OnRecordAfterDeleteRequest.Subscribe(count(c, "OnRecordAfterDeleteRequest"))
	// app.Hooks.OnRecordAuthRequest().Add(count[*core.RecordAuthEvent](c, "OnRecordAuthRequest"))
	// app.Hooks.OnRecordBeforeAuthWithPasswordRequest().Add(count[*core.RecordAuthWithPasswordEvent](c, "OnRecordBeforeAuthWithPasswordRequest"))
	// app.Hooks.OnRecordAfterAuthWithPasswordRequest().Add(count[*core.RecordAuthWithPasswordEvent](c, "OnRecordAfterAuthWithPasswordRequest"))
	// app.Hooks.OnRecordBeforeAuthWithOAuth2Request().Add(count[*core.RecordAuthWithOAuth2Event](c, "OnRecordBeforeAuthWithOAuth2Request"))
	// app.Hooks.OnRecordAfterAuthWithOAuth2Request().Add(count[*core.RecordAuthWithOAuth2Event](c, "OnRecordAfterAuthWithOAuth2Request"))
	// app.Hooks.OnRecordBeforeAuthRefreshRequest().Add(count[*core.RecordAuthRefreshEvent](c, "OnRecordBeforeAuthRefreshRequest"))
	// app.Hooks.OnRecordAfterAuthRefreshRequest().Add(count[*core.RecordAuthRefreshEvent](c, "OnRecordAfterAuthRefreshRequest"))
	// app.Hooks.OnRecordBeforeRequestPasswordResetRequest().Add(count[*core.RecordRequestPasswordResetEvent](c, "OnRecordBeforeRequestPasswordResetRequest"))
	// app.Hooks.OnRecordAfterRequestPasswordResetRequest().Add(count[*core.RecordRequestPasswordResetEvent](c, "OnRecordAfterRequestPasswordResetRequest"))
	// app.Hooks.OnRecordBeforeConfirmPasswordResetRequest().Add(count[*core.RecordConfirmPasswordResetEvent](c, "OnRecordBeforeConfirmPasswordResetRequest"))
	// app.Hooks.OnRecordAfterConfirmPasswordResetRequest().Add(count[*core.RecordConfirmPasswordResetEvent](c, "OnRecordAfterConfirmPasswordResetRequest"))
	// app.Hooks.OnRecordBeforeRequestVerificationRequest().Add(count[*core.RecordRequestVerificationEvent](c, "OnRecordBeforeRequestVerificationRequest"))
	// app.Hooks.OnRecordAfterRequestVerificationRequest().Add(count[*core.RecordRequestVerificationEvent](c, "OnRecordAfterRequestVerificationRequest"))
	// app.Hooks.OnRecordBeforeConfirmVerificationRequest().Add(count[*core.RecordConfirmVerificationEvent](c, "OnRecordBeforeConfirmVerificationRequest"))
	// app.Hooks.OnRecordAfterConfirmVerificationRequest().Add(count[*core.RecordConfirmVerificationEvent](c, "OnRecordAfterConfirmVerificationRequest"))
	// app.Hooks.OnRecordBeforeRequestEmailChangeRequest().Add(count[*core.RecordRequestEmailChangeEvent](c, "OnRecordBeforeRequestEmailChangeRequest"))
	// app.Hooks.OnRecordAfterRequestEmailChangeRequest().Add(count[*core.RecordRequestEmailChangeEvent](c, "OnRecordAfterRequestEmailChangeRequest"))
	// app.Hooks.OnRecordBeforeConfirmEmailChangeRequest().Add(count[*core.RecordConfirmEmailChangeEvent](c, "OnRecordBeforeConfirmEmailChangeRequest"))
	// app.Hooks.OnRecordAfterConfirmEmailChangeRequest().Add(count[*core.RecordConfirmEmailChangeEvent](c, "OnRecordAfterConfirmEmailChangeRequest"))
	// app.Hooks.OnRecordListExternalAuthsRequest().Add(count[*core.RecordListExternalAuthsEvent](c, "OnRecordListExternalAuthsRequest"))
	// app.Hooks.OnRecordBeforeUnlinkExternalAuthRequest().Add(count[*core.RecordUnlinkExternalAuthEvent](c, "OnRecordBeforeUnlinkExternalAuthRequest"))
	// app.Hooks.OnRecordAfterUnlinkExternalAuthRequest().Add(count[*core.RecordUnlinkExternalAuthEvent](c, "OnRecordAfterUnlinkExternalAuthRequest"))
	// app.Hooks.OnMailerBeforeAdminResetPasswordSend().Add(count[*core.MailerAdminEvent](c, "OnMailerBeforeAdminResetPasswordSend"))
	// app.Hooks.OnMailerAfterAdminResetPasswordSend().Add(count[*core.MailerAdminEvent](c, "OnMailerAfterAdminResetPasswordSend"))
	// app.Hooks.OnMailerBeforeRecordResetPasswordSend().Add(count[*core.MailerRecordEvent](c, "OnMailerBeforeRecordResetPasswordSend"))
	// app.Hooks.OnMailerAfterRecordResetPasswordSend().Add(count[*core.MailerRecordEvent](c, "OnMailerAfterRecordResetPasswordSend"))
	// app.Hooks.OnMailerBeforeRecordVerificationSend().Add(count[*core.MailerRecordEvent](c, "OnMailerBeforeRecordVerificationSend"))
	// app.Hooks.OnMailerAfterRecordVerificationSend().Add(count[*core.MailerRecordEvent](c, "OnMailerAfterRecordVerificationSend"))
	// app.Hooks.OnMailerBeforeRecordChangeEmailSend().Add(count[*core.MailerRecordEvent](c, "OnMailerBeforeRecordChangeEmailSend"))
	// app.Hooks.OnMailerAfterRecordChangeEmailSend().Add(count[*core.MailerRecordEvent](c, "OnMailerAfterRecordChangeEmailSend"))
	// app.Hooks.OnRealtimeConnectRequest().Add(count[*core.RealtimeConnectEvent](c, "OnRealtimeConnectRequest"))
	// app.Hooks.OnRealtimeDisconnectRequest().Add(count[*core.RealtimeDisconnectEvent](c, "OnRealtimeDisconnectRequest"))
	// app.Hooks.OnRealtimeBeforeMessageSend().Add(count[*core.RealtimeMessageEvent](c, "OnRealtimeBeforeMessageSend"))
	// app.Hooks.OnRealtimeAfterMessageSend().Add(count[*core.RealtimeMessageEvent](c, "OnRealtimeAfterMessageSend"))
	// app.Hooks.OnRealtimeBeforeSubscribeRequest().Add(count[*core.RealtimeSubscribeEvent](c, "OnRealtimeBeforeSubscribeRequest"))
	// app.Hooks.OnRealtimeAfterSubscribeRequest().Add(count[*core.RealtimeSubscribeEvent](c, "OnRealtimeAfterSubscribeRequest"))
	// app.Hooks.OnSettingsListRequest().Add(count[*core.SettingsListEvent](c, "OnSettingsListRequest"))
	// app.Hooks.OnSettingsBeforeUpdateRequest().Add(count[*core.SettingsUpdateEvent](c, "OnSettingsBeforeUpdateRequest"))
	// app.Hooks.OnSettingsAfterUpdateRequest().Add(count[*core.SettingsUpdateEvent](c, "OnSettingsAfterUpdateRequest"))
	// app.Hooks.OnCollectionsListRequest().Add(count[*core.CollectionsListEvent](c, "OnCollectionsListRequest"))
	// app.Hooks.OnCollectionViewRequest().Add(count[*core.CollectionViewEvent](c, "OnCollectionViewRequest"))
	// app.Hooks.OnCollectionBeforeCreateRequest().Add(count[*core.CollectionCreateEvent](c, "OnCollectionBeforeCreateRequest"))
	// app.Hooks.OnCollectionAfterCreateRequest().Add(count[*core.CollectionCreateEvent](c, "OnCollectionAfterCreateRequest"))
	// app.Hooks.OnCollectionBeforeUpdateRequest().Add(count[*core.CollectionUpdateEvent](c, "OnCollectionBeforeUpdateRequest"))
	// app.Hooks.OnCollectionAfterUpdateRequest().Add(count[*core.CollectionUpdateEvent](c, "OnCollectionAfterUpdateRequest"))
	// app.Hooks.OnCollectionBeforeDeleteRequest().Add(count[*core.CollectionDeleteEvent](c, "OnCollectionBeforeDeleteRequest"))
	// app.Hooks.OnCollectionAfterDeleteRequest().Add(count[*core.CollectionDeleteEvent](c, "OnCollectionAfterDeleteRequest"))
	// app.Hooks.OnCollectionsBeforeImportRequest().Add(count[*core.CollectionsImportEvent](c, "OnCollectionsBeforeImportRequest"))
	// app.Hooks.OnCollectionsAfterImportRequest().Add(count[*core.CollectionsImportEvent](c, "OnCollectionsAfterImportRequest"))
	// app.Hooks.OnAdminsListRequest().Add(count[*core.AdminsListEvent](c, "OnAdminsListRequest"))
	// app.Hooks.OnAdminViewRequest().Add(count[*core.AdminViewEvent](c, "OnAdminViewRequest"))
	// app.Hooks.OnAdminBeforeCreateRequest().Add(count[*core.AdminCreateEvent](c, "OnAdminBeforeCreateRequest"))
	// app.Hooks.OnAdminAfterCreateRequest().Add(count[*core.AdminCreateEvent](c, "OnAdminAfterCreateRequest"))
	// app.Hooks.OnAdminBeforeUpdateRequest().Add(count[*core.AdminUpdateEvent](c, "OnAdminBeforeUpdateRequest"))
	// app.Hooks.OnAdminAfterUpdateRequest().Add(count[*core.AdminUpdateEvent](c, "OnAdminAfterUpdateRequest"))
	// app.Hooks.OnAdminBeforeDeleteRequest().Add(count[*core.AdminDeleteEvent](c, "OnAdminBeforeDeleteRequest"))
	// app.Hooks.OnAdminAfterDeleteRequest().Add(count[*core.AdminDeleteEvent](c, "OnAdminAfterDeleteRequest"))
	// app.Hooks.OnAdminAuthRequest().Add(count[*core.AdminAuthEvent](c, "OnAdminAuthRequest"))
	// app.Hooks.OnAdminBeforeAuthWithPasswordRequest().Add(count[*core.AdminAuthWithPasswordEvent](c, "OnAdminBeforeAuthWithPasswordRequest"))
	// app.Hooks.OnAdminAfterAuthWithPasswordRequest().Add(count[*core.AdminAuthWithPasswordEvent](c, "OnAdminAfterAuthWithPasswordRequest"))
	// app.Hooks.OnAdminBeforeAuthRefreshRequest().Add(count[*core.AdminAuthRefreshEvent](c, "OnAdminBeforeAuthRefreshRequest"))
	// app.Hooks.OnAdminAfterAuthRefreshRequest().Add(count[*core.AdminAuthRefreshEvent](c, "OnAdminAfterAuthRefreshRequest"))
	// app.Hooks.OnAdminBeforeRequestPasswordResetRequest().Add(count[*core.AdminRequestPasswordResetEvent](c, "OnAdminBeforeRequestPasswordResetRequest"))
	// app.Hooks.OnAdminAfterRequestPasswordResetRequest().Add(count[*core.AdminRequestPasswordResetEvent](c, "OnAdminAfterRequestPasswordResetRequest"))
	// app.Hooks.OnAdminBeforeConfirmPasswordResetRequest().Add(count[*core.AdminConfirmPasswordResetEvent](c, "OnAdminBeforeConfirmPasswordResetRequest"))
	// app.Hooks.OnAdminAfterConfirmPasswordResetRequest().Add(count[*core.AdminConfirmPasswordResetEvent](c, "OnAdminAfterConfirmPasswordResetRequest"))
	// app.Hooks.OnFileDownloadRequest().Add(count[*core.FileDownloadEvent](c, "OnFileDownloadRequest"))
	// app.Hooks.OnFileBeforeTokenRequest().Add(count[*core.FileTokenEvent](c, "OnFileBeforeTokenRequest"))
	// app.Hooks.OnFileAfterTokenRequest().Add(count[*core.FileTokenEvent](c, "OnFileAfterTokenRequest"))
	// app.Hooks.OnFileAfterTokenRequest().Add(count[*core.FileTokenEvent](c, "OnFileAfterTokenRequest"))

	return c
}

func count(c *Counter, name string) func(ctx context.Context, table string, record any) {
	return func(_ context.Context, _ string, _ any) {
		c.Increment(name)
	}
}
