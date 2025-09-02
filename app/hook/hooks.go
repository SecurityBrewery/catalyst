package hook

type Hooks struct {
	OnRecordsListRequest        *Hook
	OnRecordViewRequest         *Hook
	OnRecordBeforeCreateRequest *Hook
	OnRecordAfterCreateRequest  *Hook
	OnRecordBeforeUpdateRequest *Hook
	OnRecordAfterUpdateRequest  *Hook
	OnRecordBeforeDeleteRequest *Hook
	OnRecordAfterDeleteRequest  *Hook
}

func NewHooks() *Hooks {
	return &Hooks{
		OnRecordsListRequest:        &Hook{},
		OnRecordViewRequest:         &Hook{},
		OnRecordBeforeCreateRequest: &Hook{},
		OnRecordAfterCreateRequest:  &Hook{},
		OnRecordBeforeUpdateRequest: &Hook{},
		OnRecordAfterUpdateRequest:  &Hook{},
		OnRecordBeforeDeleteRequest: &Hook{},
		OnRecordAfterDeleteRequest:  &Hook{},
	}
}
