package testing

import (
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

	baseApp.Settings().Logs.MaxDays = 0

	defaultTestData(t, baseApp)

	counter := countEvents(baseApp)

	return baseApp, counter, func() { _ = os.RemoveAll(temp) }
}

func generateAdminToken(t *testing.T, baseApp core.App, email string) (string, error) {
	t.Helper()

	admin, err := baseApp.Dao().FindAdminByEmail(email)
	if err != nil {
		return "", err
	}

	return tokens.NewAdminAuthToken(baseApp, admin)
}

func generateRecordToken(t *testing.T, baseApp core.App, email string) (string, error) {
	t.Helper()

	record, err := baseApp.Dao().FindAuthRecordByEmail(migrations.UserCollectionName, email)
	if err != nil {
		return "", err
	}

	return tokens.NewRecordAuthToken(baseApp, record)
}

func countEvents(t *pocketbase.PocketBase) *Counter {
	c := NewCounter()

	t.OnBeforeApiError().Add(func(_ *core.ApiErrorEvent) error {
		c.Increment("OnBeforeApiError")

		return nil
	})

	t.OnAfterApiError().Add(func(_ *core.ApiErrorEvent) error {
		c.Increment("OnAfterApiError")

		return nil
	})

	t.OnModelBeforeCreate().Add(func(_ *core.ModelEvent) error {
		c.Increment("OnModelBeforeCreate")

		return nil
	})

	t.OnModelAfterCreate().Add(func(_ *core.ModelEvent) error {
		c.Increment("OnModelAfterCreate")

		return nil
	})

	t.OnModelBeforeUpdate().Add(func(_ *core.ModelEvent) error {
		c.Increment("OnModelBeforeUpdate")

		return nil
	})

	t.OnModelAfterUpdate().Add(func(_ *core.ModelEvent) error {
		c.Increment("OnModelAfterUpdate")

		return nil
	})

	t.OnModelBeforeDelete().Add(func(_ *core.ModelEvent) error {
		c.Increment("OnModelBeforeDelete")

		return nil
	})

	t.OnModelAfterDelete().Add(func(_ *core.ModelEvent) error {
		c.Increment("OnModelAfterDelete")

		return nil
	})

	t.OnRecordsListRequest().Add(func(_ *core.RecordsListEvent) error {
		c.Increment("OnRecordsListRequest")

		return nil
	})

	t.OnRecordViewRequest().Add(func(_ *core.RecordViewEvent) error {
		c.Increment("OnRecordViewRequest")

		return nil
	})

	t.OnRecordBeforeCreateRequest().Add(func(_ *core.RecordCreateEvent) error {
		c.Increment("OnRecordBeforeCreateRequest")

		return nil
	})

	t.OnRecordAfterCreateRequest().Add(func(_ *core.RecordCreateEvent) error {
		c.Increment("OnRecordAfterCreateRequest")

		return nil
	})

	t.OnRecordBeforeUpdateRequest().Add(func(_ *core.RecordUpdateEvent) error {
		c.Increment("OnRecordBeforeUpdateRequest")

		return nil
	})

	t.OnRecordAfterUpdateRequest().Add(func(_ *core.RecordUpdateEvent) error {
		c.Increment("OnRecordAfterUpdateRequest")

		return nil
	})

	t.OnRecordBeforeDeleteRequest().Add(func(_ *core.RecordDeleteEvent) error {
		c.Increment("OnRecordBeforeDeleteRequest")

		return nil
	})

	t.OnRecordAfterDeleteRequest().Add(func(_ *core.RecordDeleteEvent) error {
		c.Increment("OnRecordAfterDeleteRequest")

		return nil
	})

	t.OnRecordAuthRequest().Add(func(_ *core.RecordAuthEvent) error {
		c.Increment("OnRecordAuthRequest")

		return nil
	})

	t.OnRecordBeforeAuthWithPasswordRequest().Add(func(_ *core.RecordAuthWithPasswordEvent) error {
		c.Increment("OnRecordBeforeAuthWithPasswordRequest")

		return nil
	})

	t.OnRecordAfterAuthWithPasswordRequest().Add(func(_ *core.RecordAuthWithPasswordEvent) error {
		c.Increment("OnRecordAfterAuthWithPasswordRequest")

		return nil
	})

	t.OnRecordBeforeAuthWithOAuth2Request().Add(func(_ *core.RecordAuthWithOAuth2Event) error {
		c.Increment("OnRecordBeforeAuthWithOAuth2Request")

		return nil
	})

	t.OnRecordAfterAuthWithOAuth2Request().Add(func(_ *core.RecordAuthWithOAuth2Event) error {
		c.Increment("OnRecordAfterAuthWithOAuth2Request")

		return nil
	})

	t.OnRecordBeforeAuthRefreshRequest().Add(func(_ *core.RecordAuthRefreshEvent) error {
		c.Increment("OnRecordBeforeAuthRefreshRequest")

		return nil
	})

	t.OnRecordAfterAuthRefreshRequest().Add(func(_ *core.RecordAuthRefreshEvent) error {
		c.Increment("OnRecordAfterAuthRefreshRequest")

		return nil
	})

	t.OnRecordBeforeRequestPasswordResetRequest().Add(func(_ *core.RecordRequestPasswordResetEvent) error {
		c.Increment("OnRecordBeforeRequestPasswordResetRequest")

		return nil
	})

	t.OnRecordAfterRequestPasswordResetRequest().Add(func(_ *core.RecordRequestPasswordResetEvent) error {
		c.Increment("OnRecordAfterRequestPasswordResetRequest")

		return nil
	})

	t.OnRecordBeforeConfirmPasswordResetRequest().Add(func(_ *core.RecordConfirmPasswordResetEvent) error {
		c.Increment("OnRecordBeforeConfirmPasswordResetRequest")

		return nil
	})

	t.OnRecordAfterConfirmPasswordResetRequest().Add(func(_ *core.RecordConfirmPasswordResetEvent) error {
		c.Increment("OnRecordAfterConfirmPasswordResetRequest")

		return nil
	})

	t.OnRecordBeforeRequestVerificationRequest().Add(func(_ *core.RecordRequestVerificationEvent) error {
		c.Increment("OnRecordBeforeRequestVerificationRequest")

		return nil
	})

	t.OnRecordAfterRequestVerificationRequest().Add(func(_ *core.RecordRequestVerificationEvent) error {
		c.Increment("OnRecordAfterRequestVerificationRequest")

		return nil
	})

	t.OnRecordBeforeConfirmVerificationRequest().Add(func(_ *core.RecordConfirmVerificationEvent) error {
		c.Increment("OnRecordBeforeConfirmVerificationRequest")

		return nil
	})

	t.OnRecordAfterConfirmVerificationRequest().Add(func(_ *core.RecordConfirmVerificationEvent) error {
		c.Increment("OnRecordAfterConfirmVerificationRequest")

		return nil
	})

	t.OnRecordBeforeRequestEmailChangeRequest().Add(func(_ *core.RecordRequestEmailChangeEvent) error {
		c.Increment("OnRecordBeforeRequestEmailChangeRequest")

		return nil
	})

	t.OnRecordAfterRequestEmailChangeRequest().Add(func(_ *core.RecordRequestEmailChangeEvent) error {
		c.Increment("OnRecordAfterRequestEmailChangeRequest")

		return nil
	})

	t.OnRecordBeforeConfirmEmailChangeRequest().Add(func(_ *core.RecordConfirmEmailChangeEvent) error {
		c.Increment("OnRecordBeforeConfirmEmailChangeRequest")

		return nil
	})

	t.OnRecordAfterConfirmEmailChangeRequest().Add(func(_ *core.RecordConfirmEmailChangeEvent) error {
		c.Increment("OnRecordAfterConfirmEmailChangeRequest")

		return nil
	})

	t.OnRecordListExternalAuthsRequest().Add(func(_ *core.RecordListExternalAuthsEvent) error {
		c.Increment("OnRecordListExternalAuthsRequest")

		return nil
	})

	t.OnRecordBeforeUnlinkExternalAuthRequest().Add(func(_ *core.RecordUnlinkExternalAuthEvent) error {
		c.Increment("OnRecordBeforeUnlinkExternalAuthRequest")

		return nil
	})

	t.OnRecordAfterUnlinkExternalAuthRequest().Add(func(_ *core.RecordUnlinkExternalAuthEvent) error {
		c.Increment("OnRecordAfterUnlinkExternalAuthRequest")

		return nil
	})

	t.OnMailerBeforeAdminResetPasswordSend().Add(func(_ *core.MailerAdminEvent) error {
		c.Increment("OnMailerBeforeAdminResetPasswordSend")

		return nil
	})

	t.OnMailerAfterAdminResetPasswordSend().Add(func(_ *core.MailerAdminEvent) error {
		c.Increment("OnMailerAfterAdminResetPasswordSend")

		return nil
	})

	t.OnMailerBeforeRecordResetPasswordSend().Add(func(_ *core.MailerRecordEvent) error {
		c.Increment("OnMailerBeforeRecordResetPasswordSend")

		return nil
	})

	t.OnMailerAfterRecordResetPasswordSend().Add(func(_ *core.MailerRecordEvent) error {
		c.Increment("OnMailerAfterRecordResetPasswordSend")

		return nil
	})

	t.OnMailerBeforeRecordVerificationSend().Add(func(_ *core.MailerRecordEvent) error {
		c.Increment("OnMailerBeforeRecordVerificationSend")

		return nil
	})

	t.OnMailerAfterRecordVerificationSend().Add(func(_ *core.MailerRecordEvent) error {
		c.Increment("OnMailerAfterRecordVerificationSend")

		return nil
	})

	t.OnMailerBeforeRecordChangeEmailSend().Add(func(_ *core.MailerRecordEvent) error {
		c.Increment("OnMailerBeforeRecordChangeEmailSend")

		return nil
	})

	t.OnMailerAfterRecordChangeEmailSend().Add(func(_ *core.MailerRecordEvent) error {
		c.Increment("OnMailerAfterRecordChangeEmailSend")

		return nil
	})

	t.OnRealtimeConnectRequest().Add(func(_ *core.RealtimeConnectEvent) error {
		c.Increment("OnRealtimeConnectRequest")

		return nil
	})

	t.OnRealtimeDisconnectRequest().Add(func(_ *core.RealtimeDisconnectEvent) error {
		c.Increment("OnRealtimeDisconnectRequest")

		return nil
	})

	t.OnRealtimeBeforeMessageSend().Add(func(_ *core.RealtimeMessageEvent) error {
		c.Increment("OnRealtimeBeforeMessageSend")

		return nil
	})

	t.OnRealtimeAfterMessageSend().Add(func(_ *core.RealtimeMessageEvent) error {
		c.Increment("OnRealtimeAfterMessageSend")

		return nil
	})

	t.OnRealtimeBeforeSubscribeRequest().Add(func(_ *core.RealtimeSubscribeEvent) error {
		c.Increment("OnRealtimeBeforeSubscribeRequest")

		return nil
	})

	t.OnRealtimeAfterSubscribeRequest().Add(func(_ *core.RealtimeSubscribeEvent) error {
		c.Increment("OnRealtimeAfterSubscribeRequest")

		return nil
	})

	t.OnSettingsListRequest().Add(func(_ *core.SettingsListEvent) error {
		c.Increment("OnSettingsListRequest")

		return nil
	})

	t.OnSettingsBeforeUpdateRequest().Add(func(_ *core.SettingsUpdateEvent) error {
		c.Increment("OnSettingsBeforeUpdateRequest")

		return nil
	})

	t.OnSettingsAfterUpdateRequest().Add(func(_ *core.SettingsUpdateEvent) error {
		c.Increment("OnSettingsAfterUpdateRequest")

		return nil
	})

	t.OnCollectionsListRequest().Add(func(_ *core.CollectionsListEvent) error {
		c.Increment("OnCollectionsListRequest")

		return nil
	})

	t.OnCollectionViewRequest().Add(func(_ *core.CollectionViewEvent) error {
		c.Increment("OnCollectionViewRequest")

		return nil
	})

	t.OnCollectionBeforeCreateRequest().Add(func(_ *core.CollectionCreateEvent) error {
		c.Increment("OnCollectionBeforeCreateRequest")

		return nil
	})

	t.OnCollectionAfterCreateRequest().Add(func(_ *core.CollectionCreateEvent) error {
		c.Increment("OnCollectionAfterCreateRequest")

		return nil
	})

	t.OnCollectionBeforeUpdateRequest().Add(func(_ *core.CollectionUpdateEvent) error {
		c.Increment("OnCollectionBeforeUpdateRequest")

		return nil
	})

	t.OnCollectionAfterUpdateRequest().Add(func(_ *core.CollectionUpdateEvent) error {
		c.Increment("OnCollectionAfterUpdateRequest")

		return nil
	})

	t.OnCollectionBeforeDeleteRequest().Add(func(_ *core.CollectionDeleteEvent) error {
		c.Increment("OnCollectionBeforeDeleteRequest")

		return nil
	})

	t.OnCollectionAfterDeleteRequest().Add(func(_ *core.CollectionDeleteEvent) error {
		c.Increment("OnCollectionAfterDeleteRequest")

		return nil
	})

	t.OnCollectionsBeforeImportRequest().Add(func(_ *core.CollectionsImportEvent) error {
		c.Increment("OnCollectionsBeforeImportRequest")

		return nil
	})

	t.OnCollectionsAfterImportRequest().Add(func(_ *core.CollectionsImportEvent) error {
		c.Increment("OnCollectionsAfterImportRequest")

		return nil
	})

	t.OnAdminsListRequest().Add(func(_ *core.AdminsListEvent) error {
		c.Increment("OnAdminsListRequest")

		return nil
	})

	t.OnAdminViewRequest().Add(func(_ *core.AdminViewEvent) error {
		c.Increment("OnAdminViewRequest")

		return nil
	})

	t.OnAdminBeforeCreateRequest().Add(func(_ *core.AdminCreateEvent) error {
		c.Increment("OnAdminBeforeCreateRequest")

		return nil
	})

	t.OnAdminAfterCreateRequest().Add(func(_ *core.AdminCreateEvent) error {
		c.Increment("OnAdminAfterCreateRequest")

		return nil
	})

	t.OnAdminBeforeUpdateRequest().Add(func(_ *core.AdminUpdateEvent) error {
		c.Increment("OnAdminBeforeUpdateRequest")

		return nil
	})

	t.OnAdminAfterUpdateRequest().Add(func(_ *core.AdminUpdateEvent) error {
		c.Increment("OnAdminAfterUpdateRequest")

		return nil
	})

	t.OnAdminBeforeDeleteRequest().Add(func(_ *core.AdminDeleteEvent) error {
		c.Increment("OnAdminBeforeDeleteRequest")

		return nil
	})

	t.OnAdminAfterDeleteRequest().Add(func(_ *core.AdminDeleteEvent) error {
		c.Increment("OnAdminAfterDeleteRequest")

		return nil
	})

	t.OnAdminAuthRequest().Add(func(_ *core.AdminAuthEvent) error {
		c.Increment("OnAdminAuthRequest")

		return nil
	})

	t.OnAdminBeforeAuthWithPasswordRequest().Add(func(_ *core.AdminAuthWithPasswordEvent) error {
		c.Increment("OnAdminBeforeAuthWithPasswordRequest")

		return nil
	})

	t.OnAdminAfterAuthWithPasswordRequest().Add(func(_ *core.AdminAuthWithPasswordEvent) error {
		c.Increment("OnAdminAfterAuthWithPasswordRequest")

		return nil
	})

	t.OnAdminBeforeAuthRefreshRequest().Add(func(_ *core.AdminAuthRefreshEvent) error {
		c.Increment("OnAdminBeforeAuthRefreshRequest")

		return nil
	})

	t.OnAdminAfterAuthRefreshRequest().Add(func(_ *core.AdminAuthRefreshEvent) error {
		c.Increment("OnAdminAfterAuthRefreshRequest")

		return nil
	})

	t.OnAdminBeforeRequestPasswordResetRequest().Add(func(_ *core.AdminRequestPasswordResetEvent) error {
		c.Increment("OnAdminBeforeRequestPasswordResetRequest")

		return nil
	})

	t.OnAdminAfterRequestPasswordResetRequest().Add(func(_ *core.AdminRequestPasswordResetEvent) error {
		c.Increment("OnAdminAfterRequestPasswordResetRequest")

		return nil
	})

	t.OnAdminBeforeConfirmPasswordResetRequest().Add(func(_ *core.AdminConfirmPasswordResetEvent) error {
		c.Increment("OnAdminBeforeConfirmPasswordResetRequest")

		return nil
	})

	t.OnAdminAfterConfirmPasswordResetRequest().Add(func(_ *core.AdminConfirmPasswordResetEvent) error {
		c.Increment("OnAdminAfterConfirmPasswordResetRequest")

		return nil
	})

	t.OnFileDownloadRequest().Add(func(_ *core.FileDownloadEvent) error {
		c.Increment("OnFileDownloadRequest")

		return nil
	})

	t.OnFileBeforeTokenRequest().Add(func(_ *core.FileTokenEvent) error {
		c.Increment("OnFileBeforeTokenRequest")

		return nil
	})

	t.OnFileAfterTokenRequest().Add(func(_ *core.FileTokenEvent) error {
		c.Increment("OnFileAfterTokenRequest")

		return nil
	})

	return c
}
