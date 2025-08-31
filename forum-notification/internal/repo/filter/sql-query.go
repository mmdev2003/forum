package filter

const GetFiltersByAccountId = `
SELECT enabled_types
		FROM notification_settings
		WHERE account_id = @account_id;
`

const UpsertFilterByAccountId = `
INSERT INTO notification_settings (account_id, enabled_types)
		VALUES (@account_id, @enabled_types::notification_type_enum[])
		ON CONFLICT (account_id) 
		DO UPDATE SET enabled_types = EXCLUDED.enabled_types
`

const CheckIfNotificationTypeIsEnabled = `
SELECT @notification_type::notification_type_enum = ANY(enabled_types)
        FROM notification_settings
        WHERE account_id = @account_id
`
