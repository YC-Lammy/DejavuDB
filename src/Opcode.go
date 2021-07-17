package main

const (
	/* Opcode:  Goto * P2 * * *
	**
	** An unconditional jump to address P2.
	** The next instruction executed will be
	** the one at index P2 from the beginning of
	** the program.
	**
	** The P1 parameter is not actually used by this opcode.  However, it
	** is sometimes set to 1 instead of 0 as a hint to the command-line shell
	** that this Goto is the bottom of a loop and that the lines from P2 down
	** to the current line should be indented for EXPLAIN output.
	 */
	Op_Goto = 0x00

	/* Opcode:  Gosub P1 P2 * * *
	**
	** Write the current address onto register P1
	** and then jump to address P2.
	 */
	Op_Gosub = 0x01

	/* Opcode:  Return P1 * * * *
	**
	** Jump to the next instruction after the address in register P1.  After
	** the jump, register P1 becomes undefined.
	 */
	Op_Return = 0x02

	/* Opcode: InitCoroutine P1 P2 P3 * *
	**
	** Set up register P1 so that it will Yield to the coroutine
	** located at address P3.
	**
	** If P2!=0 then the coroutine implementation immediately follows
	** this opcode.  So jump over the coroutine implementation to
	** address P2.
	**
	** See also: EndCoroutine
	 */
	OP_InitCoroutine = 0x03

	/* Opcode:  EndCoroutine P1 * * * *
	 **
	 ** The instruction at the address in register P1 is a Yield.
	 ** Jump to the P2 parameter of that Yield.
	 ** After the jump, register P1 becomes undefined.
	 **
	 ** See also: InitCoroutine
	 */
	OP_EndCoroutine = 0x04

	/* Opcode:  Yield P1 P2 * * *
	 **
	 ** Swap the program counter with the value in register P1.  This
	 ** has the effect of yielding to a coroutine.
	 **
	 ** If the coroutine that is launched by this instruction ends with
	 ** Yield or Return then continue to the next instruction.  But if
	 ** the coroutine launched by this instruction ends with
	 ** EndCoroutine, then jump to P2 rather than continuing with the
	 ** next instruction.
	 **
	 ** See also: InitCoroutine
	 */
	Op_Yield = 0x05

	/* Opcode:  HaltIfNull  P1 P2 P3 P4 P5
	 ** Synopsis: if r[P3]=null halt
	 **
	 ** Check the value in register P3.  If it is NULL then Halt using
	 ** parameter P1, P2, and P4 as if this were a Halt instruction.  If the
	 ** value in register P3 is not NULL, then this routine is a no-op.
	 ** The P5 parameter should be 1.
	 */
	Op_HaltIfNull = 0x06

	/* Opcode:  Halt P1 P2 * P4 P5
	 **
	 ** Exit immediately.  All open cursors, etc are closed
	 ** automatically.
	 **
	 ** P1 is the result code returned by sqlite3_exec(), sqlite3_reset(),
	 ** or sqlite3_finalize().  For a normal halt, this should be SQLITE_OK (0).
	 ** For errors, it can be some other value.  If P1!=0 then P2 will determine
	 ** whether or not to rollback the current transaction.  Do not rollback
	 ** if P2==OE_Fail. Do the rollback if P2==OE_Rollback.  If P2==OE_Abort,
	 ** then back out all changes that have occurred during this execution of the
	 ** VDBE, but do not rollback the transaction.
	 **
	 ** If P4 is not null then it is an error message string.
	 **
	 ** P5 is a value between 0 and 4, inclusive, that modifies the P4 string.
	 **
	 **    0:  (no change)
	 **    1:  NOT NULL contraint failed: P4
	 **    2:  UNIQUE constraint failed: P4
	 **    3:  CHECK constraint failed: P4
	 **    4:  FOREIGN KEY constraint failed: P4
	 **
	 ** If P5 is not zero and P4 is NULL, then everything after the ":" is
	 ** omitted.
	 **
	 ** There is an implied "Halt 0 0 0" instruction inserted at the very end of
	 ** every program.  So a jump past the last instruction of the program
	 ** is the same as executing Halt.
	 */
	Op_Halt = 0x07

	/* Opcode: Integer P1 P2 * * *
	** Synopsis: r[P2]=P1
	**
	** The 32-bit integer value P1 is written into register P2.
	 */
	OP_Integer = 0x08

	/* Opcode: Int64 * P2 * P4 *
	** Synopsis: r[P2]=P4
	**
	** P4 is a pointer to a 64-bit integer value.
	** Write that value into register P2.
	 */
	Op_Int64 = 0x09

	/* Opcode: Int64 * P2 * P4 *
	** Synopsis: r[P2]=P4
	**
	** P4 is a pointer to a 128-bit integer value.
	** Write that value into register P2.
	 */
	Op_Int128 = 0x0a

	/* Opcode: Real * P2 * P4 *
	** Synopsis: r[P2]=P4
	**
	** P4 is a pointer to a 64-bit floating point value.
	** Write that value into register P2.
	 */
	Op_Real = 0x0b

	/* Opcode: String8 * P2 * P4 *
	** Synopsis: r[P2]='P4'
	**
	** P4 points to a nul terminated UTF-8 string. This opcode is transformed
	** into a String opcode before it is executed for the first time.  During
	** this transformation, the length of string P4 is computed and stored
	** as the P1 parameter.
	 */
	Op_String8 = 0x0c

	/* Opcode: String P1 P2 P3 P4 P5
	** Synopsis: r[P2]='P4' (len=P1)
	**
	** The string value P4 of length P1 (bytes) is stored in register P2.
	**
	** If P3 is not zero and the content of register P3 is equal to P5, then
	** the datatype of the register P2 is converted to BLOB.  The content is
	** the same sequence of bytes, it is merely interpreted as a BLOB instead
	** of a string, as if it had been CAST.  In other words:
	**
	** if( P3!=0 and reg[P3]==P5 ) reg[P2] := CAST(reg[P2] as BLOB)
	 */
	Op_String = 0x0d

	/* Opcode: Null P1 P2 P3 * *
	** Synopsis: r[P2..P3]=NULL
	**
	** Write a NULL into registers P2.  If P3 greater than P2, then also write
	** NULL into register P3 and every register in between P2 and P3.  If P3
	** is less than P2 (typically P3 is zero) then only register P2 is
	** set to NULL.
	**
	** If the P1 value is non-zero, then also set the MEM_Cleared flag so that
	** NULL values will not compare equal even if SQLITE_NULLEQ is set on
	** OP_Ne or OP_Eq.
	 */
	Op_Null = 0x0e

	/* Opcode: SoftNull P1 * * * *
	** Synopsis: r[P1]=NULL
	**
	** Set register P1 to have the value NULL as seen by the OP_MakeRecord
	** instruction, but do not free any string or blob memory associated with
	** the register, so that if the value was a string or blob that was
	** previously copied using OP_SCopy, the copies will continue to be valid.
	 */
	Op_SoftNull = 0x0f

	/* Opcode: Blob P1 P2 * P4 *
	** Synopsis: r[P2]=P4 (len=P1)
	**
	** P4 points to a blob of data P1 bytes long.  Store this
	** blob in register P2.
	 */
	Op_Blob = 0x10

	/* Opcode: Variable P1 P2 * P4 *
	 ** Synopsis: r[P2]=parameter(P1,P4)
	 **
	 ** Transfer the values of bound parameter P1 into register P2
	 **
	 ** If the parameter is named, then its name appears in P4.
	 ** The P4 value is used by sqlite3_bind_parameter_name().
	 */
	Op_Variable = 0x11

	/* Opcode: Move P1 P2 P3 * *
	 ** Synopsis: r[P2@P3]=r[P1@P3]
	 **
	 ** Move the P3 values in register P1..P1+P3-1 over into
	 ** registers P2..P2+P3-1.  Registers P1..P1+P3-1 are
	 ** left holding a NULL.  It is an error for register ranges
	 ** P1..P1+P3-1 and P2..P2+P3-1 to overlap.  It is an error
	 ** for P3 to be less than 1.
	 */
	Op_Move = 0x12

	/* Opcode: Copy P1 P2 P3 * *
	** Synopsis: r[P2@P3+1]=r[P1@P3+1]
	**
	** Make a copy of registers P1..P1+P3 into registers P2..P2+P3.
	**
	** This instruction makes a deep copy of the value.  A duplicate
	** is made of any string or blob constant.  See also OP_SCopy.
	 */
	Op_Copy = 0x13

	/* Opcode: SCopy P1 P2 * * *
	 ** Synopsis: r[P2]=r[P1]
	 **
	 ** Make a shallow copy of register P1 into register P2.
	 **
	 ** This instruction makes a shallow copy of the value.  If the value
	 ** is a string or blob, then the copy is only a pointer to the
	 ** original and hence if the original changes so will the copy.
	 ** Worse, if the original is deallocated, the copy becomes invalid.
	 ** Thus the program must guarantee that the original will not change
	 ** during the lifetime of the copy.  Use OP_Copy to make a complete
	 ** copy.
	 */
	Op_SCopy = 0x14

	/* Opcode: IntCopy P1 P2 * * *
	 ** Synopsis: r[P2]=r[P1]
	 **
	 ** Transfer the integer value held in register P1 into register P2.
	 **
	 ** This is an optimized version of SCopy that works only for integer
	 ** values.
	 */
	Op_IntCopy = 0x15

	/* Opcode: ChngCntRow P1 P2 * * *
	** Synopsis: output=r[P1]
	**
	** Output value in register P1 as the chance count for a DML statement,
	** due to the "PRAGMA count_changes=ON" setting.  Or, if there was a
	** foreign key error in the statement, trigger the error now.
	**
	** This opcode is a variant of OP_ResultRow that checks the foreign key
	** immediate constraint count and throws an error if the count is
	** non-zero.  The P2 opcode must be 1.
	 */
	Op_ChngCntRow = 0x16

	/* Opcode: ResultRow P1 P2 * * *
	 ** Synopsis: output=r[P1@P2]
	 **
	 ** The registers P1 through P1+P2-1 contain a single row of
	 ** results. This opcode causes the sqlite3_step() call to terminate
	 ** with an SQLITE_ROW return code and it sets up the sqlite3_stmt
	 ** structure to provide access to the r(P1)..r(P1+P2-1) values as
	 ** the result row.
	 */
	OP_ResultRow = 0x17

	/* Opcode: Concat P1 P2 P3 * *
	** Synopsis: r[P3]=r[P2]+r[P1]
	**
	** Add the text in register P1 onto the end of the text in
	** register P2 and store the result in register P3.
	** If either the P1 or P2 text are NULL then store NULL in P3.
	**
	**   P3 = P2 || P1
	**
	** It is illegal for P1 and P3 to be the same register. Sometimes,
	** if P3 is the same register as P2, the implementation is able
	** to avoid a memcpy().
	 */
	Op_Concat = 0x18

	/* Opcode: Add P1 P2 P3 * *
	** Synopsis: r[P3]=r[P1]+r[P2]
	**
	** Add the value in register P1 to the value in register P2
	** and store the result in register P3.
	** If either input is NULL, the result is NULL.
	 */
	Op_Add = 0x19

	/* Opcode: Multiply P1 P2 P3 * *
	 ** Synopsis: r[P3]=r[P1]*r[P2]
	 **
	 **
	 ** Multiply the value in register P1 by the value in register P2
	 ** and store the result in register P3.
	 ** If either input is NULL, the result is NULL.
	 */
	Op_Multiply = 0x1a

	/* Opcode: Subtract P1 P2 P3 * *
	 ** Synopsis: r[P3]=r[P2]-r[P1]
	 **
	 ** Subtract the value in register P1 from the value in register P2
	 ** and store the result in register P3.
	 ** If either input is NULL, the result is NULL.
	 */
	Op_Subtract = 0x1b

	/* Opcode: Divide P1 P2 P3 * *
	 ** Synopsis: r[P3]=r[P2]/r[P1]
	 **
	 ** Divide the value in register P1 by the value in register P2
	 ** and store the result in register P3 (P3=P2/P1). If the value in
	 ** register P1 is zero, then the result is NULL. If either input is
	 ** NULL, the result is NULL.
	 */
	Op_Divide = 0x1c

	/* Opcode: Remainder P1 P2 P3 * *
	 ** Synopsis: r[P3]=r[P2]%r[P1]
	 **
	 ** Compute the remainder after integer register P2 is divided by
	 ** register P1 and store the result in register P3.
	 ** If the value in register P1 is zero the result is NULL.
	 ** If either operand is NULL, the result is NULL.
	 */
	Op_Reminder = 0x1d

	/* Opcode: CollSeq P1 * * P4
	**
	** P4 is a pointer to a CollSeq object. If the next call to a user function
	** or aggregate calls sqlite3GetFuncCollSeq(), this collation sequence will
	** be returned. This is used by the built-in min(), max() and nullif()
	** functions.
	**
	** If P1 is not zero, then it is a register that a subsequent min() or
	** max() aggregate will set to 1 if the current row is not the minimum or
	** maximum.  The P1 register is initialized to 0 by this instruction.
	**
	** The interface used by the implementation of the aforementioned functions
	** to retrieve the collation sequence set by this opcode is not available
	** publicly.  Only built-in functions have access to this feature.
	 */
	Op_CollSeq = 0x1e

	/* Opcode: BitAnd P1 P2 P3 * *
	 ** Synopsis: r[P3]=r[P1]&r[P2]
	 **
	 ** Take the bit-wise AND of the values in register P1 and P2 and
	 ** store the result in register P3.
	 ** If either input is NULL, the result is NULL.
	 */
	Op_BitAnd = 0x1f

	/* Opcode: BitOr P1 P2 P3 * *
	 ** Synopsis: r[P3]=r[P1]|r[P2]
	 **
	 ** Take the bit-wise OR of the values in register P1 and P2 and
	 ** store the result in register P3.
	 ** If either input is NULL, the result is NULL.
	 */
	Op_BitOr = 0x20

	/* Opcode: ShiftLeft P1 P2 P3 * *
	 ** Synopsis: r[P3]=r[P2]<<r[P1]
	 **
	 ** Shift the integer value in register P2 to the left by the
	 ** number of bits specified by the integer in register P1.
	 ** Store the result in register P3.
	 ** If either input is NULL, the result is NULL.
	 */
	Op_ShiftLeft = 0x21

	/* Opcode: ShiftRight P1 P2 P3 * *
	 ** Synopsis: r[P3]=r[P2]>>r[P1]
	 **
	 ** Shift the integer value in register P2 to the right by the
	 ** number of bits specified by the integer in register P1.
	 ** Store the result in register P3.
	 ** If either input is NULL, the result is NULL.
	 */
	Op_ShiftRight = 0x22

	/* Opcode: AddImm  P1 P2 * * *
	 ** Synopsis: r[P1]=r[P1]+P2
	 **
	 ** Add the constant P2 to the value in register P1.
	 ** The result is always an integer.
	 **
	 ** To force any register to be an integer, just add 0.
	 */
	Op_AddImm = 0x23

	/* Opcode: MustBeInt P1 P2 * * *
	**
	** Force the value in register P1 to be an integer.  If the value
	** in P1 is not an integer and cannot be converted into an integer
	** without data loss, then jump immediately to P2, or if P2==0
	** raise an SQLITE_MISMATCH exception.
	 */
	Op_MustBeInt = 0x24
)
