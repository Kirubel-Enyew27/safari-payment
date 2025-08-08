-- name: SavePayment :one
INSERT INTO payments (
    checkout_request_id,
    merchant_request_id,
    phone_number,
    amount,
    mpesa_receipt,
    transaction_date,
    result_code,
    result_desc
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetPaymentByCheckoutRequestID :one
SELECT * FROM payments
WHERE checkout_request_id = $1;

-- name: ListPayments :many
SELECT * FROM payments
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;