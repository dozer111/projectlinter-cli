package common_sets

import "github.com/dozer111/projectlinter-core/rules"

// leaf - this is commonly used practise for sets in projectlinter
//
// # The main reason is - readability
//
// As for me - it is easier to read the code like
// NewRuleTree
//
//	leaf
//	leaf
//		leaf
//		leaf
//	leaf
//	leaf
//		leaf
//	leaf
//
// instead of
//
// NewRuleTree
//
//	rules.NewLeaf
//	rules.NewLeaf
//		rules.NewLeaf
//		rules.NewLeaf
//	rules.NewLeaf
//
// ...
func leaf(r rules.Rule, children ...rules.RuleTreeLeaf) rules.RuleTreeLeaf {
	return rules.NewLeaf(r, children...)
}
