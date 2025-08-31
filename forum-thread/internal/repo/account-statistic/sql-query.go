package account_statistic

var CreateAccountStatistic = `
INSERT INTO account_statistics (account_id)
VALUES (@account_id)
RETURNING id;
`

var AddSentMessagesToTopicsCount = `
UPDATE account_statistics
SET sent_messages_to_topics_count = sent_messages_to_topics_count + 1
WHERE account_id = @account_id;
`

var AddCreatedTopicsCount = `
UPDATE account_statistics
SET created_topics_count = created_topics_count + 1
`

var StatisticByAccountID = `
SELECT * FROM account_statistics
WHERE account_id = @account_id;
`
