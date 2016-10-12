- 单人匹配
```
sequenceDiagram
Gate->>Lobby: C2L_JoinQueue
Lobby->>Session: PlayerState:L2S_JoinQueue
Session->>Lobby: S2L_JoinQueue
Lobby->>Gate: L2C_JoinQueue
```