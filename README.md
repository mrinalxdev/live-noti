## Currently building this project

Workflow 

                 +---------------------------+
                 |       Client Browser      |
                 +---------------------------+
                            |
                            | 1. WebSocket Request (/ws)
                            v
                 +---------------------------+
                 |        HTTP Server        |
                 |  (Gorilla WebSocket)      |
                 +---------------------------+
                            |
                            | 2. Upgrade HTTP to WebSocket
                            v
                 +---------------------------+
                 |     WebSocket Connection  |
                 |       (Handle Messages)   |
                 +---------------------------+
                            |
      +---------------------+---------------------+
      |                                           |
3. join_channel                              4. user_action
      |                                           |
      v                                           v
+------------------+                      +------------------+
|  Add Connection  |                      |  Publish Action  |
|  to Clients Map  |                      |  to Redis PubSub |
+------------------+                      +------------------+
      |                                           |
      v                                           v
+------------------+                      +------------------+
|  Add User to     |                      | Redis Notifies   |
|  Redis Channel   |                      | Subscribers      |
+------------------+                      +------------------+
      |                                           |
      v                                           v
+------------------+                      +------------------+
| Subscribe to     |                      | Broadcast Action |
| Redis Channel    |                      | to WebSocket     |
| in Goroutine     |                      | Clients          |
+------------------+                      +------------------+
      |                                           |
      v                                           v
+------------------+                      +------------------+
| Receive Messages |                      | WebSocket Clients|
| from Redis PubSub|                      | Receive Message  |
| (Goroutine)      |                      | from Server      |
+------------------+                      +------------------+

Might be using more technologies to integrate .