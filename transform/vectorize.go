package transform

import (
	"strings"

	"tinygo.org/x/go-llvm"
)

// VectorizePass looks for calls to vectorize(...) and attempts to replace them with SIMD operations.
func VectorizePass(mod llvm.Module) bool {
	modified := false

	ctx := llvm.GlobalContext()

	for fn := mod.FirstFunction(); !fn.IsNil(); fn = llvm.NextFunction(fn) {
		for bb := fn.FirstBasicBlock(); !bb.IsNil(); bb = llvm.NextBasicBlock(bb) {
			for inst := bb.FirstInstruction(); !inst.IsNil(); inst = llvm.NextInstruction(inst) {
				if !inst.IsACallInst().IsNil() {
					callee := inst.CalledValue()
					name := callee.Name()
					if strings.Contains(name, "vectorize") {
						// Example: replace with dummy <4 x i32> add pattern
						builder := ctx.NewBuilder()
						builder.SetInsertPointBefore(inst)

						int32Ty := ctx.Int32Type()
						vecTy := llvm.VectorType(int32Ty, 4)

						src := builder.CreateAlloca(vecTy, "src")
						dst := builder.CreateAlloca(vecTy, "dst")

						load := builder.CreateLoad(vecTy, src, "v")
						ones := llvm.ConstVector([]llvm.Value{
							llvm.ConstInt(int32Ty, 1, false),
							llvm.ConstInt(int32Ty, 1, false),
							llvm.ConstInt(int32Ty, 1, false),
							llvm.ConstInt(int32Ty, 1, false),
						}, false)
						add := builder.CreateAdd(load, ones, "vadd")
						builder.CreateStore(add, dst)

						inst.EraseFromParentAsInstruction()
						modified = true
					}
				}
			}
		}
	}

	return modified
}

// func lower() {
// 	vecTy := llvm.VectorType(llvm.Int32Type(), 4) // <4 x i32>
// 	v := builder.CreateLoad(vecTy, srcPtr, "")
// 	ones := llvm.ConstVector([]llvm.Value{
// 		llvm.ConstInt(llvm.Int32Type(), 1, false),
// 		llvm.ConstInt(llvm.Int32Type(), 1, false),
// 		llvm.ConstInt(llvm.Int32Type(), 1, false),
// 		llvm.ConstInt(llvm.Int32Type(), 1, false),
// 	})
// 	res := builder.CreateAdd(v, ones, "")
// 	builder.CreateStore(res, dstPtr)

// }
