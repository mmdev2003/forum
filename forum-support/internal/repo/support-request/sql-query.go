package support_request

const CreateSupportRequest = `
INSERT INTO support_requests (account_id, title, description)
		VALUES (@account_id, @title, @description)
RETURNING id;
`

const SetSupportRequest = `
UPDATE support_requests
SET 
    status = @status,
    updated_at = CURRENT_TIMESTAMP
WHERE id = @request_id;
`

const GetRequestById = `
SELECT id, account_id, title, description, status, created_at, updated_at FROM support_requests
    WHERE id = @request_id;
`

const GetRequests = `
    SELECT * FROM support_requests
`

const GetRequestsWithStatus = GetRequests + `
    WHERE status = @status::support_request_status
`
