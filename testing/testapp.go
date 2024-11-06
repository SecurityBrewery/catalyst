package testing

import (
	"fmt"
	"os"
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tokens"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/migrations"
)

func App(t *testing.T) (*pocketbase.PocketBase, *Counter, func()) {
	t.Helper()

	temp, err := os.MkdirTemp("", "catalyst_test_data")
	if err != nil {
		t.Fatal(err)
	}

	baseApp, err := app.App(temp, true)
	if err != nil {
		t.Fatal(err)
	}

	if err := baseApp.Bootstrap(); err != nil {
		t.Fatal(fmt.Errorf("failed to bootstrap: %w", err))
	}

	baseApp.Settings().Logs.MaxDays = 0

	defaultTestData(t, baseApp)

	counter := countEvents(baseApp)

	return baseApp, counter, func() { _ = os.RemoveAll(temp) }
}

func generateAdminToken(t *testing.T, baseApp core.App, email string) (string, error) {
	t.Helper()

	admin, err := baseApp.Dao().FindAdminByEmail(email)
	if err != nil {
		return "", fmt.Errorf("failed to find admin: %w", err)
	}

	return tokens.NewAdminAuthToken(baseApp, admin)
}

func generateRecordToken(t *testing.T, baseApp core.App, email string) (string, error) {
	t.Helper()

	record, err := baseApp.Dao().FindAuthRecordByEmail(migrations.UserCollectionName, email)
	if err != nil {
		return "", fmt.Errorf("failed to find record: %w", err)
	}

	return tokens.NewRecordAuthToken(baseApp, record)
}

func countEvents(t *pocketbase.PocketBase) *Counter {
	c := NewCounter()

	t.OnBeforeApiError().Add(count[*core.ApiErrorEvent](c, "OnBeforeApiError"))
	t.OnBeforeApiError().Add(count[*core.ApiErrorEvent](c, "OnBeforeApiError"))
	t.OnAfterApiError().Add(count[*core.ApiErrorEvent](c, "OnAfterApiError"))
	t.OnModelBeforeCreate().Add(count[*core.ModelEvent](c, "OnModelBeforeCreate"))
	t.OnModelAfterCreate().Add(count[*core.ModelEvent](c, "OnModelAfterCreate"))
	t.OnModelBeforeUpdate().Add(count[*core.ModelEvent](c, "OnModelBeforeUpdate"))
	t.OnModelAfterUpdate().Add(count[*core.ModelEvent](c, "OnModelAfterUpdate"))
	t.OnModelBeforeDelete().Add(count[*core.ModelEvent](c, "OnModelBeforeDelete"))
	t.OnModelAfterDelete().Add(count[*core.ModelEvent](c, "OnModelAfterDelete"))
	t.OnRecordsListRequest().Add(count[*core.RecordsListEvent](c, "OnRecordsListRequest"))
	t.OnRecordViewRequest().Add(count[*core.RecordViewEvent](c, "OnRecordViewRequest"))
	t.OnRecordBeforeCreateRequest().Add(count[*core.RecordCreateEvent](c, "OnRecordBeforeCreateRequest"))
	t.OnRecordAfterCreateRequest().Add(count[*core.RecordCreateEvent](c, "OnRecordAfterCreateRequest"))
	t.OnRecordBeforeUpdateRequest().Add(count[*core.RecordUpdateEvent](c, "OnRecordBeforeUpdateRequest"))
	t.OnRecordAfterUpdateRequest().Add(count[*core.RecordUpdateEvent](c, "OnRecordAfterUpdateRequest"))
	t.OnRecordBeforeDeleteRequest().Add(count[*core.RecordDeleteEvent](c, "OnRecordBeforeDeleteRequest"))
	t.OnRecordAfterDeleteRequest().Add(count[*core.RecordDeleteEvent](c, "OnRecordAfterDeleteRequest"))
	t.OnRecordAuthRequest().Add(count[*core.RecordAuthEvent](c, "OnRecordAuthRequest"))
	t.OnRecordBeforeAuthWithPasswordRequest().Add(count[*core.RecordAuthWithPasswordEvent](c, "OnRecordBeforeAuthWithPasswordRequest"))
	t.OnRecordAfterAuthWithPasswordRequest().Add(count[*core.RecordAuthWithPasswordEvent](c, "OnRecordAfterAuthWithPasswordRequest"))
	t.OnRecordBeforeAuthWithOAuth2Request().Add(count[*core.RecordAuthWithOAuth2Event](c, "OnRecordBeforeAuthWithOAuth2Request"))
	t.OnRecordAfterAuthWithOAuth2Request().Add(count[*core.RecordAuthWithOAuth2Event](c, "OnRecordAfterAuthWithOAuth2Request"))
	t.OnRecordBeforeAuthRefreshRequest().Add(count[*core.RecordAuthRefreshEvent](c, "OnRecordBeforeAuthRefreshRequest"))
	t.OnRecordAfterAuthRefreshRequest().Add(count[*core.RecordAuthRefreshEvent](c, "OnRecordAfterAuthRefreshRequest"))
	t.OnRecordBeforeRequestPasswordResetRequest().Add(count[*core.RecordRequestPasswordResetEvent](c, "OnRecordBeforeRequestPasswordResetRequest"))
	t.OnRecordAfterRequestPasswordResetRequest().Add(count[*core.RecordRequestPasswordResetEvent](c, "OnRecordAfterRequestPasswordResetRequest"))
	t.OnRecordBeforeConfirmPasswordResetRequest().Add(count[*core.RecordConfirmPasswordResetEvent](c, "OnRecordBeforeConfirmPasswordResetRequest"))
	t.OnRecordAfterConfirmPasswordResetRequest().Add(count[*core.RecordConfirmPasswordResetEvent](c, "OnRecordAfterConfirmPasswordResetRequest"))
	t.OnRecordBeforeRequestVerificationRequest().Add(count[*core.RecordRequestVerificationEvent](c, "OnRecordBeforeRequestVerificationRequest"))
	t.OnRecordAfterRequestVerificationRequest().Add(count[*core.RecordRequestVerificationEvent](c, "OnRecordAfterRequestVerificationRequest"))
	t.OnRecordBeforeConfirmVerificationRequest().Add(count[*core.RecordConfirmVerificationEvent](c, "OnRecordBeforeConfirmVerificationRequest"))
	t.OnRecordAfterConfirmVerificationRequest().Add(count[*core.RecordConfirmVerificationEvent](c, "OnRecordAfterConfirmVerificationRequest"))
	t.OnRecordBeforeRequestEmailChangeRequest().Add(count[*core.RecordRequestEmailChangeEvent](c, "OnRecordBeforeRequestEmailChangeRequest"))
	t.OnRecordAfterRequestEmailChangeRequest().Add(count[*core.RecordRequestEmailChangeEvent](c, "OnRecordAfterRequestEmailChangeRequest"))
	t.OnRecordBeforeConfirmEmailChangeRequest().Add(count[*core.RecordConfirmEmailChangeEvent](c, "OnRecordBeforeConfirmEmailChangeRequest"))
	t.OnRecordAfterConfirmEmailChangeRequest().Add(count[*core.RecordConfirmEmailChangeEvent](c, "OnRecordAfterConfirmEmailChangeRequest"))
	t.OnRecordListExternalAuthsRequest().Add(count[*core.RecordListExternalAuthsEvent](c, "OnRecordListExternalAuthsRequest"))
	t.OnRecordBeforeUnlinkExternalAuthRequest().Add(count[*core.RecordUnlinkExternalAuthEvent](c, "OnRecordBeforeUnlinkExternalAuthRequest"))
	t.OnRecordAfterUnlinkExternalAuthRequest().Add(count[*core.RecordUnlinkExternalAuthEvent](c, "OnRecordAfterUnlinkExternalAuthRequest"))
	t.OnMailerBeforeAdminResetPasswordSend().Add(count[*core.MailerAdminEvent](c, "OnMailerBeforeAdminResetPasswordSend"))
	t.OnMailerAfterAdminResetPasswordSend().Add(count[*core.MailerAdminEvent](c, "OnMailerAfterAdminResetPasswordSend"))
	t.OnMailerBeforeRecordResetPasswordSend().Add(count[*core.MailerRecordEvent](c, "OnMailerBeforeRecordResetPasswordSend"))
	t.OnMailerAfterRecordResetPasswordSend().Add(count[*core.MailerRecordEvent](c, "OnMailerAfterRecordResetPasswordSend"))
	t.OnMailerBeforeRecordVerificationSend().Add(count[*core.MailerRecordEvent](c, "OnMailerBeforeRecordVerificationSend"))
	t.OnMailerAfterRecordVerificationSend().Add(count[*core.MailerRecordEvent](c, "OnMailerAfterRecordVerificationSend"))
	t.OnMailerBeforeRecordChangeEmailSend().Add(count[*core.MailerRecordEvent](c, "OnMailerBeforeRecordChangeEmailSend"))
	t.OnMailerAfterRecordChangeEmailSend().Add(count[*core.MailerRecordEvent](c, "OnMailerAfterRecordChangeEmailSend"))
	t.OnRealtimeConnectRequest().Add(count[*core.RealtimeConnectEvent](c, "OnRealtimeConnectRequest"))
	t.OnRealtimeDisconnectRequest().Add(count[*core.RealtimeDisconnectEvent](c, "OnRealtimeDisconnectRequest"))
	t.OnRealtimeBeforeMessageSend().Add(count[*core.RealtimeMessageEvent](c, "OnRealtimeBeforeMessageSend"))
	t.OnRealtimeAfterMessageSend().Add(count[*core.RealtimeMessageEvent](c, "OnRealtimeAfterMessageSend"))
	t.OnRealtimeBeforeSubscribeRequest().Add(count[*core.RealtimeSubscribeEvent](c, "OnRealtimeBeforeSubscribeRequest"))
	t.OnRealtimeAfterSubscribeRequest().Add(count[*core.RealtimeSubscribeEvent](c, "OnRealtimeAfterSubscribeRequest"))
	t.OnSettingsListRequest().Add(count[*core.SettingsListEvent](c, "OnSettingsListRequest"))
	t.OnSettingsBeforeUpdateRequest().Add(count[*core.SettingsUpdateEvent](c, "OnSettingsBeforeUpdateRequest"))
	t.OnSettingsAfterUpdateRequest().Add(count[*core.SettingsUpdateEvent](c, "OnSettingsAfterUpdateRequest"))
	t.OnCollectionsListRequest().Add(count[*core.CollectionsListEvent](c, "OnCollectionsListRequest"))
	t.OnCollectionViewRequest().Add(count[*core.CollectionViewEvent](c, "OnCollectionViewRequest"))
	t.OnCollectionBeforeCreateRequest().Add(count[*core.CollectionCreateEvent](c, "OnCollectionBeforeCreateRequest"))
	t.OnCollectionAfterCreateRequest().Add(count[*core.CollectionCreateEvent](c, "OnCollectionAfterCreateRequest"))
	t.OnCollectionBeforeUpdateRequest().Add(count[*core.CollectionUpdateEvent](c, "OnCollectionBeforeUpdateRequest"))
	t.OnCollectionAfterUpdateRequest().Add(count[*core.CollectionUpdateEvent](c, "OnCollectionAfterUpdateRequest"))
	t.OnCollectionBeforeDeleteRequest().Add(count[*core.CollectionDeleteEvent](c, "OnCollectionBeforeDeleteRequest"))
	t.OnCollectionAfterDeleteRequest().Add(count[*core.CollectionDeleteEvent](c, "OnCollectionAfterDeleteRequest"))
	t.OnCollectionsBeforeImportRequest().Add(count[*core.CollectionsImportEvent](c, "OnCollectionsBeforeImportRequest"))
	t.OnCollectionsAfterImportRequest().Add(count[*core.CollectionsImportEvent](c, "OnCollectionsAfterImportRequest"))
	t.OnAdminsListRequest().Add(count[*core.AdminsListEvent](c, "OnAdminsListRequest"))
	t.OnAdminViewRequest().Add(count[*core.AdminViewEvent](c, "OnAdminViewRequest"))
	t.OnAdminBeforeCreateRequest().Add(count[*core.AdminCreateEvent](c, "OnAdminBeforeCreateRequest"))
	t.OnAdminAfterCreateRequest().Add(count[*core.AdminCreateEvent](c, "OnAdminAfterCreateRequest"))
	t.OnAdminBeforeUpdateRequest().Add(count[*core.AdminUpdateEvent](c, "OnAdminBeforeUpdateRequest"))
	t.OnAdminAfterUpdateRequest().Add(count[*core.AdminUpdateEvent](c, "OnAdminAfterUpdateRequest"))
	t.OnAdminBeforeDeleteRequest().Add(count[*core.AdminDeleteEvent](c, "OnAdminBeforeDeleteRequest"))
	t.OnAdminAfterDeleteRequest().Add(count[*core.AdminDeleteEvent](c, "OnAdminAfterDeleteRequest"))
	t.OnAdminAuthRequest().Add(count[*core.AdminAuthEvent](c, "OnAdminAuthRequest"))
	t.OnAdminBeforeAuthWithPasswordRequest().Add(count[*core.AdminAuthWithPasswordEvent](c, "OnAdminBeforeAuthWithPasswordRequest"))
	t.OnAdminAfterAuthWithPasswordRequest().Add(count[*core.AdminAuthWithPasswordEvent](c, "OnAdminAfterAuthWithPasswordRequest"))
	t.OnAdminBeforeAuthRefreshRequest().Add(count[*core.AdminAuthRefreshEvent](c, "OnAdminBeforeAuthRefreshRequest"))
	t.OnAdminAfterAuthRefreshRequest().Add(count[*core.AdminAuthRefreshEvent](c, "OnAdminAfterAuthRefreshRequest"))
	t.OnAdminBeforeRequestPasswordResetRequest().Add(count[*core.AdminRequestPasswordResetEvent](c, "OnAdminBeforeRequestPasswordResetRequest"))
	t.OnAdminAfterRequestPasswordResetRequest().Add(count[*core.AdminRequestPasswordResetEvent](c, "OnAdminAfterRequestPasswordResetRequest"))
	t.OnAdminBeforeConfirmPasswordResetRequest().Add(count[*core.AdminConfirmPasswordResetEvent](c, "OnAdminBeforeConfirmPasswordResetRequest"))
	t.OnAdminAfterConfirmPasswordResetRequest().Add(count[*core.AdminConfirmPasswordResetEvent](c, "OnAdminAfterConfirmPasswordResetRequest"))
	t.OnFileDownloadRequest().Add(count[*core.FileDownloadEvent](c, "OnFileDownloadRequest"))
	t.OnFileBeforeTokenRequest().Add(count[*core.FileTokenEvent](c, "OnFileBeforeTokenRequest"))
	t.OnFileAfterTokenRequest().Add(count[*core.FileTokenEvent](c, "OnFileAfterTokenRequest"))
	t.OnFileAfterTokenRequest().Add(count[*core.FileTokenEvent](c, "OnFileAfterTokenRequest"))

	return c
}

func count[T any](c *Counter, name string) func(_ T) error {
	return func(_ T) error {
		c.Increment(name)

		return nil
	}
}
