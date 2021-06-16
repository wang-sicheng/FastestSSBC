package chain

import (
	"encoding/hex"
	"github.com/cloudflare/cfssl/log"
	"github.com/fastestssbc/merkle"
	"github.com/fastestssbc/meta"
	"github.com/fastestssbc/util"
	"time"
)

var (
	BlockChain   = make([]meta.Block, 0)
	CurrentBlock meta.Block
	NewBlock     meta.Block
)

func init() {
	//首先创建创世区块
	GenerateGenesisBlock()
}

//创建创世区块
func GenerateGenesisBlock() {
	genesisBlock := meta.Block{}
	genesisBlock.Hash = util.CalBlockHash(genesisBlock)
	genesisBlock.Merkle = GenerateMerkleRoot(genesisBlock)
	//log.Info("GenesisBlock: ", genesisBlock)
	BlockChain = append(BlockChain, genesisBlock)
	CurrentBlock = genesisBlock
	log.Info("Block Init Successfully.")
}

//创建新区块
func GenerateNewBlock(oldBlock meta.Block, newBlock meta.Block) meta.Block {
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Merkle = GenerateMerkleRoot(newBlock)
	newBlock.Signature = "Signature"
	newBlock.Hash = util.CalBlockHash(newBlock)
	return newBlock
}

//保存区块
func StoreNewBlock() {
	BlockChain = append(BlockChain, NewBlock)
	CurrentBlock = NewBlock
	log.Info("Store NewBlock SuccessFully!")
}

func VerifyBlock(block meta.Block) bool {
	voteBool := false
	if verifyBlock(block) {
		voteBool = true
	}
	log.Info("Verify Block: ", voteBool)
	return voteBool
}
func verifyBlock(block meta.Block) bool {
	//验证逻辑 验签 验证交易 验证merkle tree root
	if block.PrevHash != CurrentBlock.Hash {
		log.Error("区块验证失败：区块Hash值非法")
		//log.Info("block.PrevHash=",block.PrevHash)
		//log.Info("CurrentBlock.Hash=",CurrentBlock.Hash)
		return false
	}
	if block.Signature != "Signature" {
		log.Info("verify block: Signature mismatch")
		return false
	}
	return verifyBlockTx(block)
}

func verifyBlockTx(b meta.Block) bool {
	if len(b.TX) != len(TransHashDataMap) {
		log.Error("所收主节点的区块交易与本地交易数量不一致")
		return false
	}
	return true
}

//生成merkleTree
func GenerateMerkleRoot(b meta.Block) string {
	mt := merkle.NewMerkleTree(transToByte(b.TX))
	return hex.EncodeToString(mt.RootNode.Data)
}

func transToByte(trans []meta.Transaction) [][]byte {
	res := [][]byte{}
	for _, data := range trans {
		res = append(res, transTobyte(data))
	}
	return res
}
func transTobyte(tran meta.Transaction) []byte {
	tranString := tran.From + tran.To + tran.Timestamp + tran.Signature + tran.Message
	return []byte(tranString)
}
