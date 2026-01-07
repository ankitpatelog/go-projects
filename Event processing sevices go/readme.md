FLOW 1: Incoming Event (Happy Path)
Example
POST /events

STEP-BY-STEP
1️⃣ Client sends webhook
   |
2️⃣ Gin Router receives request
   |
3️⃣ Rate Limiter Middleware (Redis)
   |
4️⃣ Request Validation
   |
5️⃣ Idempotency Check (Redis)
   |
6️⃣ Save Event in MySQL (PENDING)
   |
7️⃣ Push Event to Queue (Channel)
   |
8️⃣ Return 200 Accepted
   |
9️⃣ Worker picks event
   |
🔟 Event Handler executes
   |
1️⃣1️⃣ Update MySQL → PROCESSED
   |
1️⃣2️⃣ Mark Redis Idempotency Key

🔴 FLOW 2: Duplicate Event (Idempotency Flow)
Scenario

Webhook provider retries same event.

POST /events (same event_id)
   |
Rate Limiter ✔
   |
Validation ✔
   |
Redis Idempotency EXISTS ❌
   |
Return 200 OK (ignore)

Important

❌ No DB write

❌ No queue

❌ No worker

👉 Duplicate avoided

🟡 FLOW 3: Rate Limit Exceeded
POST /events
   |
Redis INCR(rate:ip)
   |
Count > limit
   |
Return 429 Too Many Requests


👉 Protects system from abuse

🔵 FLOW 4: Worker Processing (Success)
Worker picks event
   |
Update MySQL → PROCESSING
   |
Call handler (by event_type)
   |
Handler success
   |
Update MySQL → PROCESSED
   |
Set Redis Idempotency Key

🔴 FLOW 5: Worker Processing (Failure + Retry)
Worker picks event
   |
Handler fails
   |
Retry count < max?
   |
YES
   |
Increment retry_count (MySQL)
   |
Push back to queue

⚫ FLOW 6: Max Retry Reached
Worker picks event
   |
Handler fails
   |
Retry count >= max
   |
Update MySQL → FAILED
   |
Stop retrying


👉 Prevents infinite loops

🟣 FLOW 7: System Restart Recovery
App restarts
   |
Load PENDING / RETRY events from MySQL
   |
Push to queue
   |
Workers resume processing


👉 Redis lost? No issue.
👉 MySQL = source of truth

🧠 FLOW 8: Internal Code-Level Flow
main.go
 |
 |-- setup config (MySQL, Redis)
 |-- init queue
 |-- start workers
 |-- start Gin server
 |
 |--> /events
       |
       middleware/rate_limiter.go
       |
       api/event_handler.go
       |
       services/event_service.go
       |
       repository/event_repo.go
       |
       queue/channel.go
       |
       workers/worker.go
       |
       handlers/*

🧱 FLOW 9: Redis Key Lifecycle
Rate Limiting
rate:{ip}
TTL = 1 min

Idempotency
idem:event:{event_id}
TTL = 1 hour

🧠 FLOW 10: Status Lifecycle (MySQL)
PENDING → PROCESSING → PROCESSED
             |
             ↓
           RETRY → FAILED