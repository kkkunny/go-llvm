//===- llvm/Analysis/MemoryProfileInfo.h - memory profile info ---*- C++ -*-==//
//
// Part of the LLVM Project, under the Apache License v2.0 with LLVM Exceptions.
// See https://llvm.org/LICENSE.txt for license information.
// SPDX-License-Identifier: Apache-2.0 WITH LLVM-exception
//
//===----------------------------------------------------------------------===//
//
// This file contains utilities to analyze memory profile information.
//
//===----------------------------------------------------------------------===//

#ifndef LLVM_ANALYSIS_MEMORYPROFILEINFO_H
#define LLVM_ANALYSIS_MEMORYPROFILEINFO_H

#include "llvm/IR/Constants.h"
#include "llvm/IR/InstrTypes.h"
#include "llvm/IR/Metadata.h"
#include "llvm/IR/Module.h"
#include <map>

namespace llvm {
namespace memprof {

// Allocation type assigned to an allocation reached by a given context.
// More can be added but initially this is just noncold and cold.
// Values should be powers of two so that they can be ORed, in particular to
// track allocations that have different behavior with different calling
// contexts.
enum class AllocationType : uint8_t { None = 0, NotCold = 1, Cold = 2 };

/// Return the allocation type for a given set of memory profile values.
AllocationType getAllocType(uint64_t MaxAccessCount, uint64_t MinSize,
                            uint64_t MinLifetime);

/// Build callstack metadata from the provided list of call stack ids. Returns
/// the resulting metadata node.
MDNode *buildCallstackMetadata(ArrayRef<uint64_t> CallStack, LLVMContext &Ctx);

/// Returns the stack node from an MIB metadata node.
MDNode *getMIBStackNode(const MDNode *MIB);

/// Returns the allocation type from an MIB metadata node.
AllocationType getMIBAllocType(const MDNode *MIB);

/// Class to build a trie of call stack contexts for a particular profiled
/// allocation call, along with their associated allocation types.
/// The allocation will be at the root of the trie, which is then used to
/// compute the minimum lists of context ids needed to associate a call context
/// with a single allocation type.
class CallStackTrie {
private:
  struct CallStackTrieNode {
    // Allocation types for call context sharing the context prefix at this
    // node.
    uint8_t AllocTypes;
    // Map of caller stack id to the corresponding child Trie node.
    std::map<uint64_t, CallStackTrieNode *> Callers;
    CallStackTrieNode(AllocationType Type)
        : AllocTypes(static_cast<uint8_t>(Type)) {}
  };

  // The node for the allocation at the root.
  CallStackTrieNode *Alloc;
  // The allocation's leaf stack id.
  uint64_t AllocStackId;

  void deleteTrieNode(CallStackTrieNode *Node) {
    if (!Node)
      return;
    for (auto C : Node->Callers)
      deleteTrieNode(C.second);
    delete Node;
  }

  // Recursive helper to trim contexts and create metadata nodes.
  bool buildMIBNodes(CallStackTrieNode *Node, LLVMContext &Ctx,
                     std::vector<uint64_t> &MIBCallStack,
                     std::vector<Metadata *> &MIBNodes,
                     bool CalleeHasAmbiguousCallerContext);

public:
  CallStackTrie() : Alloc(nullptr), AllocStackId(0) {}
  ~CallStackTrie() { deleteTrieNode(Alloc); }

  bool empty() const { return Alloc == nullptr; }

  /// Add a call stack context with the given allocation type to the Trie.
  /// The context is represented by the list of stack ids (computed during
  /// matching via a debug location hash), expected to be in order from the
  /// allocation call down to the bottom of the call stack (i.e. callee to
  /// caller order).
  void addCallStack(AllocationType AllocType, ArrayRef<uint64_t> StackIds);

  /// Add the call stack context along with its allocation type from the MIB
  /// metadata to the Trie.
  void addCallStack(MDNode *MIB);

  /// Build and attach the minimal necessary MIB metadata. If the alloc has a
  /// single allocation type, add a function attribute instead. The reason for
  /// adding an attribute in this case is that it matches how the behavior for
  /// allocation calls will be communicated to lib call simplification after
  /// cloning or another optimization to distinguish the allocation types,
  /// which is lower overhead and more direct than maintaining this metadata.
  /// Returns true if memprof metadata attached, false if not (attribute added).
  bool buildAndAttachMIBMetadata(CallBase *CI);
};

} // end namespace memprof
} // end namespace llvm

#endif
