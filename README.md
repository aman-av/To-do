🧠 Mapping gRPC Patterns to TODO Project
✅ Unary RPC
Use case: Get a TODO item by ID

Request: GetTodoRequest { id }

Response: Todo

🔁 Server-Side Streaming
Use case: List all TODOs

Client sends filter request once

Server streams each matching TODO

🔄 Client-Side Streaming
Use case: Bulk create TODOs

Client streams multiple CreateTodoRequests

Server responds with a single BulkCreateResponse

🔁🔄 Bidirectional Streaming
Use case: Real-time TODO sync (for fun or learning)

Client streams actions (create, update, delete)

Server streams updates or acknowledgments in real time
