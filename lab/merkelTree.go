package lab

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type TreeNode struct {
	Value       string
	Left, Right *TreeNode
}

func hashMerkleLeafe(t Transaction) string {
	v, _ := json.Marshal(t)
	hash := sha256.New()
	hash.Write(v)
	return hex.EncodeToString(hash.Sum(nil))
}

func createMerkleNode(left, right *TreeNode) TreeNode {
	lv, rv := "", ""
	if left != nil {
		lv = left.Value
	}
	if right != nil {
		rv = right.Value
	}

	hash := sha256.New()
	hash.Write([]byte(lv + rv))
	return TreeNode{Value: hex.EncodeToString(hash.Sum(nil)), Left: left, Right: right}
}

func buildMerkleLayer(hashes []TreeNode) TreeNode {
	if len(hashes) == 1 {
		return hashes[0]
	}

	layer := make([]TreeNode, 0)
	for i := 0; i < len(hashes); i += 2 {
		var node TreeNode
		if i+1 < len(hashes) {
			node = createMerkleNode(&hashes[i], &hashes[i+1])
		} else {
			node = createMerkleNode(&hashes[i], nil)
		}
		layer = append(layer, node)
	}

	return buildMerkleLayer(layer)
}

func BuildMerkleTree(transactions []Transaction) TreeNode {
	leaves := make([]TreeNode, 0, len(transactions))
	for _, t := range transactions {
		leaves = append(leaves, TreeNode{Value: hashMerkleLeafe(t)})
	}
	return buildMerkleLayer(leaves)
}
