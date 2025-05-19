# TotalMassAssignmentCalculator 流程图

```mermaid
graph TD
    A[开始] --> B[初始化计算器]
    B --> C[填充初始总质量]
    C --> D{循环处理数据}
    
    D --> E[获取中子/桥梁下一总质量]
    E --> F{中子值 < 桥梁值?}
    
    F --> |是| G{中子初始总质量 ≥ 桥梁初始值?}
    G --> |是| H[生成'001'结果]
    G --> |否| I[计算纯度并生成'002'结果]
    
    F --> |否| J{中子值 > 桥梁值?}
    J --> |是| K{中子初始总质量 > 桥梁初始值?}
    K --> |是| L[按比例计算纯度]
    K --> |否| M[生成'002'错误结果]
    
    J --> |否| N[生成'00'结果]
    
    H & I & L & M & N --> O[更新总质量]
    O --> P[移动数据索引]
    P --> Q{是否还有数据?}
    Q --> |是| D
    Q --> |否| R[生成最终预测结果]
    
    R --> S[数据分组统计]
    S --> T[计算均值/方差]
    T --> U[记录分析日志]
    U --> V[结束]
    
    style A fill:#9f9,stroke:#333
    style V fill:#f99,stroke:#333
    style F,J,K,G fill:#fcf,stroke:#333
    style H,I,L,M,N fill:#cff,stroke:#333
