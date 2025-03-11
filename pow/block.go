package pow

import(
    "time"
    "crypto/sha256"
    "bytes"
)

var (
    DefaultBits = 24
)

// Block 表示区块链中的一个区块
type Block struct{
    // 区块头
    PrevBlockHash []byte        // 前一个区块的哈希值
    MerkleRoot    []byte        // 默克尔树根，用于快速校验交易
    Timestamp     int64         // 区块创建时间戳
    Bits         uint32         // 难度目标值
    Nonce        uint32         // 工作量证明的随机数

    // 区块体
    Transactions []*Transaction // 交易列表
}


type Blockchain struct {
	blocks []*Block
}


// Transaction 表示一个交易
type Transaction struct {
    ID   []byte    // 交易ID（交易数据的SHA256哈希）
    Data []byte    // 交易数据
}

// CalculateID 计算交易的ID
func (tx *Transaction) CalculateID() []byte {
    if tx.Data == nil {
        return nil
    }
    // 对交易数据进行SHA256哈希
    hash := sha256.Sum256(tx.Data)
    return hash[:]
}

// CreateBlock 创建新区块
func CreateBlock(prevBlockHash []byte, transactions []*Transaction) *Block {
    block := &Block{
        PrevBlockHash: prevBlockHash,
        Transactions:  transactions,
        Nonce:        0,
        Timestamp:    time.Now().Unix(),
        Bits:         DefaultBits,
    }
    
    // 确保所有交易都有ID
    for _, tx := range transactions {
        if tx.ID == nil {
            tx.ID = tx.CalculateID()
        }
    }
    
    // 计算默克尔根
    block.MerkleRoot = block.CalculateMerkleRoot()
    
    return block
}

// CalculateMerkleRoot 使用交易ID计算默克尔树根
func (b *Block) CalculateMerkleRoot() []byte {
    var txIDs [][]byte
    
    // 收集所有交易的ID
    for _, tx := range b.Transactions {
        if tx.ID == nil {
            tx.ID = tx.CalculateID()
        }
        txIDs = append(txIDs, tx.ID)
    }
    
    // 如果没有交易，返回空哈希
    if len(txIDs) == 0 {
        return make([]byte, 32)
    }
    
    // 构建默克尔树
    for len(txIDs) > 1 {
        // 如果是奇数个哈希，复制最后一个
        if len(txIDs)%2 != 0 {
            txIDs = append(txIDs, txIDs[len(txIDs)-1])
        }
        
        var newLevel [][]byte
        // 两两配对计算新的哈希
        for i := 0; i < len(txIDs); i += 2 {
            // 拼接相邻的两个交易ID
            concat := append(txIDs[i], txIDs[i+1]...)
            // 对拼接结果进行SHA256
            hash := sha256.Sum256(concat)
            newLevel = append(newLevel, hash[:])
        }
        txIDs = newLevel
    }
    
    return txIDs[0]
}

