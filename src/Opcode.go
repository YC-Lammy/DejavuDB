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

	/* Opcode: RealAffinity P1 * * * *
	 **
	 ** If register P1 holds an integer convert it to a real value.
	 **
	 ** This opcode is used when extracting information from a column that
	 ** has REAL affinity.  Such column values may still be stored as
	 ** integers, for space efficiency, but after extraction we want them
	 ** to have only a real value.
	 */
	OP_RealAffinity = 0x25

	/* Opcode: Cast P1 P2 * * *
	** Synopsis: affinity(r[P1])
	**
	** Force the value in register P1 to be the type defined by P2.
	**
	** <ul>
	** <li> P2=='A' &rarr; BLOB
	** <li> P2=='B' &rarr; TEXT
	** <li> P2=='C' &rarr; NUMERIC
	** <li> P2=='D' &rarr; INTEGER
	** <li> P2=='E' &rarr; REAL
	** </ul>
	**
	** A NULL value is not changed by this routine.  It remains NULL.
	 */
	Op_Cast = 0x26

	/* Opcode: Eq P1 P2 P3 P4 P5
	** Synopsis: IF r[P3]==r[P1]
	**
	** Compare the values in register P1 and P3.  If reg(P3)==reg(P1) then
	** jump to address P2.
	**
	** The SQLITE_AFF_MASK portion of P5 must be an affinity character -
	** SQLITE_AFF_TEXT, SQLITE_AFF_INTEGER, and so forth. An attempt is made
	** to coerce both inputs according to this affinity before the
	** comparison is made. If the SQLITE_AFF_MASK is 0x00, then numeric
	** affinity is used. Note that the affinity conversions are stored
	** back into the input registers P1 and P3.  So this opcode can cause
	** persistent changes to registers P1 and P3.
	**
	** Once any conversions have taken place, and neither value is NULL,
	** the values are compared. If both values are blobs then memcmp() is
	** used to determine the results of the comparison.  If both values
	** are text, then the appropriate collating function specified in
	** P4 is used to do the comparison.  If P4 is not specified then
	** memcmp() is used to compare text string.  If both values are
	** numeric, then a numeric comparison is used. If the two values
	** are of different types, then numbers are considered less than
	** strings and strings are considered less than blobs.
	**
	** If SQLITE_NULLEQ is set in P5 then the result of comparison is always either
	** true or false and is never NULL.  If both operands are NULL then the result
	** of comparison is true.  If either operand is NULL then the result is false.
	** If neither operand is NULL the result is the same as it would be if
	** the SQLITE_NULLEQ flag were omitted from P5.
	**
	** This opcode saves the result of comparison for use by the new
	** OP_Jump opcode.
	 */
	Op_Eq = 0x27

	/* Opcode: Ne P1 P2 P3 P4 P5
	** Synopsis: IF r[P3]!=r[P1]
	**
	** This works just like the Eq opcode except that the jump is taken if
	** the operands in registers P1 and P3 are not equal.  See the Eq opcode for
	** additional information.
	 */
	Op_Ne = 0x28

	/* Opcode: Lt P1 P2 P3 P4 P5
	** Synopsis: IF r[P3]<r[P1]
	**
	** Compare the values in register P1 and P3.  If reg(P3)<reg(P1) then
	** jump to address P2.
	**
	** If the SQLITE_JUMPIFNULL bit of P5 is set and either reg(P1) or
	** reg(P3) is NULL then the take the jump.  If the SQLITE_JUMPIFNULL
	** bit is clear then fall through if either operand is NULL.
	**
	** The SQLITE_AFF_MASK portion of P5 must be an affinity character -
	** SQLITE_AFF_TEXT, SQLITE_AFF_INTEGER, and so forth. An attempt is made
	** to coerce both inputs according to this affinity before the
	** comparison is made. If the SQLITE_AFF_MASK is 0x00, then numeric
	** affinity is used. Note that the affinity conversions are stored
	** back into the input registers P1 and P3.  So this opcode can cause
	** persistent changes to registers P1 and P3.
	**
	** Once any conversions have taken place, and neither value is NULL,
	** the values are compared. If both values are blobs then memcmp() is
	** used to determine the results of the comparison.  If both values
	** are text, then the appropriate collating function specified in
	** P4 is  used to do the comparison.  If P4 is not specified then
	** memcmp() is used to compare text string.  If both values are
	** numeric, then a numeric comparison is used. If the two values
	** are of different types, then numbers are considered less than
	** strings and strings are considered less than blobs.
	**
	** This opcode saves the result of comparison for use by the new
	** OP_Jump opcode.
	 */
	Op_Lt = 0x29

	/* Opcode: Le P1 P2 P3 P4 P5
	** Synopsis: IF r[P3]<=r[P1]
	**
	** This works just like the Lt opcode except that the jump is taken if
	** the content of register P3 is less than or equal to the content of
	** register P1.  See the Lt opcode for additional information.
	 */
	Op_Le = 0x2a

	/* Opcode: Gt P1 P2 P3 P4 P5
	** Synopsis: IF r[P3]>r[P1]
	**
	** This works just like the Lt opcode except that the jump is taken if
	** the content of register P3 is greater than the content of
	** register P1.  See the Lt opcode for additional information.
	 */
	Op_Gt = 0x2b

	/* Opcode: Ge P1 P2 P3 P4 P5
	** Synopsis: IF r[P3]>=r[P1]
	**
	** This works just like the Lt opcode except that the jump is taken if
	** the content of register P3 is greater than or equal to the content of
	** register P1.  See the Lt opcode for additional information.
	 */
	Op_Ge = 0x2c

	/* Opcode: ElseEq * P2 * * *
	 **
	 ** This opcode must follow an OP_Lt or OP_Gt comparison operator.  There
	 ** can be zero or more OP_ReleaseReg opcodes intervening, but no other
	 ** opcodes are allowed to occur between this instruction and the previous
	 ** OP_Lt or OP_Gt.
	 **
	 ** If result of an OP_Eq comparison on the same two operands as the
	 ** prior OP_Lt or OP_Gt would have been true, then jump to P2.
	 ** If the result of an OP_Eq comparison on the two previous
	 ** operands would have been false or NULL, then fall through.
	 */
	Op_ElseEq = 0x2d

	/* Opcode: Permutation * * * P4 *
	 **
	 ** Set the permutation used by the OP_Compare operator in the next
	 ** instruction.  The permutation is stored in the P4 operand.
	 **
	 ** The permutation is only valid until the next OP_Compare that has
	 ** the OPFLAG_PERMUTE bit set in P5. Typically the OP_Permutation should
	 ** occur immediately prior to the OP_Compare.
	 **
	 ** The first integer in the P4 integer array is the length of the array
	 ** and does not become part of the permutation.
	 */
	Op_Permutation = 0x2e

	/* Opcode: Compare P1 P2 P3 P4 P5
	** Synopsis: r[P1@P3] <-> r[P2@P3]
	**
	** Compare two vectors of registers in reg(P1)..reg(P1+P3-1) (call this
	** vector "A") and in reg(P2)..reg(P2+P3-1) ("B").  Save the result of
	** the comparison for use by the next OP_Jump instruct.
	**
	** If P5 has the OPFLAG_PERMUTE bit set, then the order of comparison is
	** determined by the most recent OP_Permutation operator.  If the
	** OPFLAG_PERMUTE bit is clear, then register are compared in sequential
	** order.
	**
	** P4 is a KeyInfo structure that defines collating sequences and sort
	** orders for the comparison.  The permutation applies to registers
	** only.  The KeyInfo elements are used sequentially.
	**
	** The comparison is a sort comparison, so NULLs compare equal,
	** NULLs are less than numbers, numbers are less than strings,
	** and strings are less than blobs.
	 */
	Op_Compare = 0x2f

	/* Opcode: Jump P1 P2 P3 * *
	**
	** Jump to the instruction at address P1, P2, or P3 depending on whether
	** in the most recent OP_Compare instruction the P1 vector was less than
	** equal to, or greater than the P2 vector, respectively.
	 */
	Op_Jump = 0x30

	/* Opcode: And P1 P2 P3 * *
	** Synopsis: r[P3]=(r[P1] && r[P2])
	**
	** Take the logical AND of the values in registers P1 and P2 and
	** write the result into register P3.
	**
	** If either P1 or P2 is 0 (false) then the result is 0 even if
	** the other input is NULL.  A NULL and true or two NULLs give
	** a NULL output.
	 */
	Op_And = 0x31

	/* Opcode: Or P1 P2 P3 * *
	** Synopsis: r[P3]=(r[P1] || r[P2])
	**
	** Take the logical OR of the values in register P1 and P2 and
	** store the answer in register P3.
	**
	** If either P1 or P2 is nonzero (true) then the result is 1 (true)
	** even if the other input is NULL.  A NULL and false or two NULLs
	** give a NULL output.
	 */
	Op_Or = 0x32

	/* Opcode: IsTrue P1 P2 P3 P4 *
	** Synopsis: r[P2] = coalesce(r[P1]==TRUE,P3) ^ P4
	**
	** This opcode implements the IS TRUE, IS FALSE, IS NOT TRUE, and
	** IS NOT FALSE operators.
	**
	** Interpret the value in register P1 as a boolean value.  Store that
	** boolean (a 0 or 1) in register P2.  Or if the value in register P1 is
	** NULL, then the P3 is stored in register P2.  Invert the answer if P4
	** is 1.
	**
	** The logic is summarized like this:
	**
	** <ul>
	** <li> If P3==0 and P4==0  then  r[P2] := r[P1] IS TRUE
	** <li> If P3==1 and P4==1  then  r[P2] := r[P1] IS FALSE
	** <li> If P3==0 and P4==1  then  r[P2] := r[P1] IS NOT TRUE
	** <li> If P3==1 and P4==0  then  r[P2] := r[P1] IS NOT FALSE
	** </ul>
	 */
	Op_IsTrue = 0x33

	/* Opcode: Not P1 P2 * * *
	** Synopsis: r[P2]= !r[P1]
	**
	** Interpret the value in register P1 as a boolean value.  Store the
	** boolean complement in register P2.  If the value in register P1 is
	** NULL, then a NULL is stored in P2.
	 */
	Op_Not = 0x34

	/* Opcode: BitNot P1 P2 * * *
	** Synopsis: r[P2]= ~r[P1]
	**
	** Interpret the content of register P1 as an integer.  Store the
	** ones-complement of the P1 value into register P2.  If P1 holds
	** a NULL then store a NULL in P2.
	 */
	Op_BitNot = 0x35

	/* Opcode: Once P1 P2 * * *
	 **
	 ** Fall through to the next instruction the first time this opcode is
	 ** encountered on each invocation of the byte-code program.  Jump to P2
	 ** on the second and all subsequent encounters during the same invocation.
	 **
	 ** Top-level programs determine first invocation by comparing the P1
	 ** operand against the P1 operand on the OP_Init opcode at the beginning
	 ** of the program.  If the P1 values differ, then fall through and make
	 ** the P1 of this opcode equal to the P1 of OP_Init.  If P1 values are
	 ** the same then take the jump.
	 **
	 ** For subprograms, there is a bitmask in the VdbeFrame that determines
	 ** whether or not the jump should be taken.  The bitmask is necessary
	 ** because the self-altering code trick does not work for recursive
	 ** triggers.
	 */
	Op_Once = 0x36

	/* Opcode: If P1 P2 P3 * *
	 **
	 ** Jump to P2 if the value in register P1 is true.  The value
	 ** is considered true if it is numeric and non-zero.  If the value
	 ** in P1 is NULL then take the jump if and only if P3 is non-zero.
	 */
	Op_If = 0x37

	/* Opcode: IfNot P1 P2 P3 * *
	**
	** Jump to P2 if the value in register P1 is False.  The value
	** is considered false if it has a numeric value of zero.  If the value
	** in P1 is NULL then take the jump if and only if P3 is non-zero.
	 */
	Op_IfNot = 0x38

	/* Opcode: IsNull P1 P2 * * *
	** Synopsis: if r[P1]==NULL goto P2
	**
	** Jump to P2 if the value in register P1 is NULL.
	 */
	Op_IsNull = 0x39

	/* Opcode: ZeroOrNull P1 P2 P3 * *
	** Synopsis: r[P2] = 0 OR NULL
	**
	** If all both registers P1 and P3 are NOT NULL, then store a zero in
	** register P2.  If either registers P1 or P3 are NULL then put
	** a NULL in register P2.
	 */
	Op_ZeroOrNull = 0x3a

	/* Opcode: NotNull P1 P2 * * *
	** Synopsis: if r[P1]!=NULL goto P2
	**
	** Jump to P2 if the value in register P1 is not NULL.
	 */
	Op_NotNull = 0x3b

	/* Opcode: IfNullRow P1 P2 P3 * *
	 ** Synopsis: if P1.nullRow then r[P3]=NULL, goto P2
	 **
	 ** Check the cursor P1 to see if it is currently pointing at a NULL row.
	 ** If it is, then set register P3 to NULL and jump immediately to P2.
	 ** If P1 is not on a NULL row, then fall through without making any
	 ** changes.
	 */
	Op_IfNullRow = 0x3c

	/* Opcode: Offset P1 P2 P3 * *
	** Synopsis: r[P3] = sqlite_offset(P1)
	**
	** Store in register r[P3] the byte offset into the database file that is the
	** start of the payload for the record at which that cursor P1 is currently
	** pointing.
	**
	** P2 is the column number for the argument to the sqlite_offset() function.
	** This opcode does not use P2 itself, but the P2 value is used by the
	** code generator.  The P1, P2, and P3 operands to this opcode are the
	** same as for OP_Column.
	**
	** This opcode is only available if SQLite is compiled with the
	** -DSQLITE_ENABLE_OFFSET_SQL_FUNC option.
	 */
	Op_Offset = 0x3d

	/* Opcode: Column P1 P2 P3 P4 P5
	** Synopsis: r[P3]=PX
	**
	** Interpret the data that cursor P1 points to as a structure built using
	** the MakeRecord instruction.  (See the MakeRecord opcode for additional
	** information about the format of the data.)  Extract the P2-th column
	** from this record.  If there are less that (P2+1)
	** values in the record, extract a NULL.
	**
	** The value extracted is stored in register P3.
	**
	** If the record contains fewer than P2 fields, then extract a NULL.  Or,
	** if the P4 argument is a P4_MEM use the value of the P4 argument as
	** the result.
	**
	** If the OPFLAG_LENGTHARG and OPFLAG_TYPEOFARG bits are set on P5 then
	** the result is guaranteed to only be used as the argument of a length()
	** or typeof() function, respectively.  The loading of large blobs can be
	** skipped for length() and all content loading can be skipped for typeof().
	 */
	Op_Column = 0x3e

	/* Opcode: Affinity P1 P2 * P4 *
	** Synopsis: affinity(r[P1@P2])
	**
	** Apply affinities to a range of P2 registers starting with P1.
	**
	** P4 is a string that is P2 characters long. The N-th character of the
	** string indicates the column affinity that should be used for the N-th
	** memory cell in the range.
	 */
	Op_Affinity = 0x3f

	/* Opcode: MakeRecord P1 P2 P3 P4 *
	 ** Synopsis: r[P3]=mkrec(r[P1@P2])
	 **
	 ** Convert P2 registers beginning with P1 into the [record format]
	 ** use as a data record in a database table or as a key
	 ** in an index.  The OP_Column opcode can decode the record later.
	 **
	 ** P4 may be a string that is P2 characters long.  The N-th character of the
	 ** string indicates the column affinity that should be used for the N-th
	 ** field of the index key.
	 **
	 ** The mapping from character to affinity is given by the SQLITE_AFF_
	 ** macros defined in sqliteInt.h.
	 **
	 ** If P4 is NULL then all index fields have the affinity BLOB.
	 **
	 ** The meaning of P5 depends on whether or not the SQLITE_ENABLE_NULL_TRIM
	 ** compile-time option is enabled:
	 **
	 **   * If SQLITE_ENABLE_NULL_TRIM is enabled, then the P5 is the index
	 **     of the right-most table that can be null-trimmed.
	 **
	 **   * If SQLITE_ENABLE_NULL_TRIM is omitted, then P5 has the value
	 **     OPFLAG_NOCHNG_MAGIC if the OP_MakeRecord opcode is allowed to
	 **     accept no-change records with serial_type 10.  This value is
	 **     only used inside an assert() and does not affect the end result.
	 */
	Op_MakeRecord = 0x40

	/* Opcode: Count P1 P2 p3 * *
	** Synopsis: r[P2]=count()
	**
	** Store the number of entries (an integer value) in the table or index
	** opened by cursor P1 in register P2.
	**
	** If P3==0, then an exact count is obtained, which involves visiting
	** every btree page of the table.  But if P3 is non-zero, an estimate
	** is returned based on the current cursor position.
	 */
	Op_Count = 0x41

	/* Opcode: Savepoint P1 * * P4 *
	**
	** Open, release or rollback the savepoint named by parameter P4, depending
	** on the value of P1. To open a new savepoint set P1==0 (SAVEPOINT_BEGIN).
	** To release (commit) an existing savepoint set P1==1 (SAVEPOINT_RELEASE).
	** To rollback an existing savepoint set P1==2 (SAVEPOINT_ROLLBACK).
	 */
	Op_Savepoint = 0x42

	/* Opcode: AutoCommit P1 P2 * * *
	 **
	 ** Set the database auto-commit flag to P1 (1 or 0). If P2 is true, roll
	 ** back any currently active btree transactions. If there are any active
	 ** VMs (apart from this one), then a ROLLBACK fails.  A COMMIT fails if
	 ** there are active writing VMs or active VMs that use shared cache.
	 **
	 ** This instruction causes the VM to halt.
	 */
	Op_AutoCommit = 0x43

	/* Opcode: Transaction P1 P2 P3 P4 P5
	**
	** Begin a transaction on database P1 if a transaction is not already
	** active.
	** If P2 is non-zero, then a write-transaction is started, or if a
	** read-transaction is already active, it is upgraded to a write-transaction.
	** If P2 is zero, then a read-transaction is started.  If P2 is 2 or more
	** then an exclusive transaction is started.
	**
	** P1 is the index of the database file on which the transaction is
	** started.  Index 0 is the main database file and index 1 is the
	** file used for temporary tables.  Indices of 2 or more are used for
	** attached databases.
	**
	** If a write-transaction is started and the Vdbe.usesStmtJournal flag is
	** true (this flag is set if the Vdbe may modify more than one row and may
	** throw an ABORT exception), a statement transaction may also be opened.
	** More specifically, a statement transaction is opened iff the database
	** connection is currently not in autocommit mode, or if there are other
	** active statements. A statement transaction allows the changes made by this
	** VDBE to be rolled back after an error without having to roll back the
	** entire transaction. If no error is encountered, the statement transaction
	** will automatically commit when the VDBE halts.
	**
	** If P5!=0 then this opcode also checks the schema cookie against P3
	** and the schema generation counter against P4.
	** The cookie changes its value whenever the database schema changes.
	** This operation is used to detect when that the cookie has changed
	** and that the current process needs to reread the schema.  If the schema
	** cookie in P3 differs from the schema cookie in the database header or
	** if the schema generation counter in P4 differs from the current
	** generation counter, then an SQLITE_SCHEMA error is raised and execution
	** halts.  The sqlite3_step() wrapper function might then reprepare the
	** statement and rerun it from the beginning.
	 */
	Op_Transaction = 0x44

	/* Opcode: ReadCookie P1 P2 P3 * *
	 **
	 ** Read cookie number P3 from database P1 and write it into register P2.
	 ** P3==1 is the schema version.  P3==2 is the database format.
	 ** P3==3 is the recommended pager cache size, and so forth.  P1==0 is
	 ** the main database file and P1==1 is the database file used to store
	 ** temporary tables.
	 **
	 ** There must be a read-lock on the database (either a transaction
	 ** must be started or there must be an open cursor) before
	 ** executing this instruction.
	 */
	Op_ReadCookie = 0x45

	/* Opcode: SetCookie P1 P2 P3 * P5
	**
	** Write the integer value P3 into cookie number P2 of database P1.
	** P2==1 is the schema version.  P2==2 is the database format.
	** P2==3 is the recommended pager cache
	** size, and so forth.  P1==0 is the main database file and P1==1 is the
	** database file used to store temporary tables.
	**
	** A transaction must be started before executing this opcode.
	**
	** If P2 is the SCHEMA_VERSION cookie (cookie number 1) then the internal
	** schema version is set to P3-P5.  The "PRAGMA schema_version=N" statement
	** has P5 set to 1, so that the internal schema version will be different
	** from the database schema version, resulting in a schema reset.
	 */
	Op_SetCookie = 0x46

	/* Opcode: OpenRead P1 P2 P3 P4 P5
	 ** Synopsis: root=P2 iDb=P3
	 **
	 ** Open a read-only cursor for the database table whose root page is
	 ** P2 in a database file.  The database file is determined by P3.
	 ** P3==0 means the main database, P3==1 means the database used for
	 ** temporary tables, and P3>1 means used the corresponding attached
	 ** database.  Give the new cursor an identifier of P1.  The P1
	 ** values need not be contiguous but all P1 values should be small integers.
	 ** It is an error for P1 to be negative.
	 **
	 ** Allowed P5 bits:
	 ** <ul>
	 ** <li>  <b>0x02 OPFLAG_SEEKEQ</b>: This cursor will only be used for
	 **       equality lookups (implemented as a pair of opcodes OP_SeekGE/OP_IdxGT
	 **       of OP_SeekLE/OP_IdxLT)
	 ** </ul>
	 **
	 ** The P4 value may be either an integer (P4_INT32) or a pointer to
	 ** a KeyInfo structure (P4_KEYINFO). If it is a pointer to a KeyInfo
	 ** object, then table being opened must be an [index b-tree] where the
	 ** KeyInfo object defines the content and collating
	 ** sequence of that index b-tree. Otherwise, if P4 is an integer
	 ** value, then the table being opened must be a [table b-tree] with a
	 ** number of columns no less than the value of P4.
	 **
	 ** See also: OpenWrite, ReopenIdx
	 */
	Op_OpenRead = 0x47

	/* Opcode: ReopenIdx P1 P2 P3 P4 P5
	** Synopsis: root=P2 iDb=P3
	**
	** The ReopenIdx opcode works like OP_OpenRead except that it first
	** checks to see if the cursor on P1 is already open on the same
	** b-tree and if it is this opcode becomes a no-op.  In other words,
	** if the cursor is already open, do not reopen it.
	**
	** The ReopenIdx opcode may only be used with P5==0 or P5==OPFLAG_SEEKEQ
	** and with P4 being a P4_KEYINFO object.  Furthermore, the P3 value must
	** be the same as every other ReopenIdx or OpenRead for the same cursor
	** number.
	**
	** Allowed P5 bits:
	** <ul>
	** <li>  <b>0x02 OPFLAG_SEEKEQ</b>: This cursor will only be used for
	**       equality lookups (implemented as a pair of opcodes OP_SeekGE/OP_IdxGT
	**       of OP_SeekLE/OP_IdxLT)
	** </ul>
	**
	** See also: OP_OpenRead, OP_OpenWrite
	 */
	Op_ReopenIdx = 0x48

	/* Opcode: OpenWrite P1 P2 P3 P4 P5
	** Synopsis: root=P2 iDb=P3
	**
	** Open a read/write cursor named P1 on the table or index whose root
	** page is P2 (or whose root page is held in register P2 if the
	** OPFLAG_P2ISREG bit is set in P5 - see below).
	**
	** The P4 value may be either an integer (P4_INT32) or a pointer to
	** a KeyInfo structure (P4_KEYINFO). If it is a pointer to a KeyInfo
	** object, then table being opened must be an [index b-tree] where the
	** KeyInfo object defines the content and collating
	** sequence of that index b-tree. Otherwise, if P4 is an integer
	** value, then the table being opened must be a [table b-tree] with a
	** number of columns no less than the value of P4.
	**
	** Allowed P5 bits:
	** <ul>
	** <li>  <b>0x02 OPFLAG_SEEKEQ</b>: This cursor will only be used for
	**       equality lookups (implemented as a pair of opcodes OP_SeekGE/OP_IdxGT
	**       of OP_SeekLE/OP_IdxLT)
	** <li>  <b>0x08 OPFLAG_FORDELETE</b>: This cursor is used only to seek
	**       and subsequently delete entries in an index btree.  This is a
	**       hint to the storage engine that the storage engine is allowed to
	**       ignore.  The hint is not used by the official SQLite b*tree storage
	**       engine, but is used by COMDB2.
	** <li>  <b>0x10 OPFLAG_P2ISREG</b>: Use the content of register P2
	**       as the root page, not the value of P2 itself.
	** </ul>
	**
	** This instruction works like OpenRead except that it opens the cursor
	** in read/write mode.
	**
	** See also: OP_OpenRead, OP_ReopenIdx
	 */
	Op_OpenWrite = 0x49

	/* Opcode: OpenDup P1 P2 * * *
	 **
	 ** Open a new cursor P1 that points to the same ephemeral table as
	 ** cursor P2.  The P2 cursor must have been opened by a prior OP_OpenEphemeral
	 ** opcode.  Only ephemeral cursors may be duplicated.
	 **
	 ** Duplicate ephemeral cursors are used for self-joins of materialized views.
	 */
	Op_OpenDup = 0x4a

	/* Opcode: OpenEphemeral P1 P2 P3 P4 P5
	 ** Synopsis: nColumn=P2
	 **
	 ** Open a new cursor P1 to a transient table.
	 ** The cursor is always opened read/write even if
	 ** the main database is read-only.  The ephemeral
	 ** table is deleted automatically when the cursor is closed.
	 **
	 ** If the cursor P1 is already opened on an ephemeral table, the table
	 ** is cleared (all content is erased).
	 **
	 ** P2 is the number of columns in the ephemeral table.
	 ** The cursor points to a BTree table if P4==0 and to a BTree index
	 ** if P4 is not 0.  If P4 is not NULL, it points to a KeyInfo structure
	 ** that defines the format of keys in the index.
	 **
	 ** The P5 parameter can be a mask of the BTREE_* flags defined
	 ** in btree.h.  These flags control aspects of the operation of
	 ** the btree.  The BTREE_OMIT_JOURNAL and BTREE_SINGLE flags are
	 ** added automatically.
	 **
	 ** If P3 is positive, then reg[P3] is modified slightly so that it
	 ** can be used as zero-length data for OP_Insert.  This is an optimization
	 ** that avoids an extra OP_Blob opcode to initialize that register.
	 */
	Op_OpenEphemeral = 0x4b

	/* Opcode: OpenAutoindex P1 P2 * P4 *
	 ** Synopsis: nColumn=P2
	 **
	 ** This opcode works the same as OP_OpenEphemeral.  It has a
	 ** different name to distinguish its use.  Tables created using
	 ** by this opcode will be used for automatically created transient
	 ** indices in joins.
	 */
	Op_OpenAutoindex = 0x4c

	/* Opcode: SorterOpen P1 P2 P3 P4 *
	 **
	 ** This opcode works like OP_OpenEphemeral except that it opens
	 ** a transient index that is specifically designed to sort large
	 ** tables using an external merge-sort algorithm.
	 **
	 ** If argument P3 is non-zero, then it indicates that the sorter may
	 ** assume that a stable sort considering the first P3 fields of each
	 ** key is sufficient to produce the required results.
	 */
	Op_SorterOpen = 0x4d

	/* Opcode: SequenceTest P1 P2 * * *
	** Synopsis: if( cursor[P1].ctr++ ) pc = P2
	**
	** P1 is a sorter cursor. If the sequence counter is currently zero, jump
	** to P2. Regardless of whether or not the jump is taken, increment the
	** the sequence value.
	 */
	Op_SequenceTest = 0x4e

	/* Opcode: OpenPseudo P1 P2 P3 * *
	 ** Synopsis: P3 columns in r[P2]
	 **
	 ** Open a new cursor that points to a fake table that contains a single
	 ** row of data.  The content of that one row is the content of memory
	 ** register P2.  In other words, cursor P1 becomes an alias for the
	 ** MEM_Blob content contained in register P2.
	 **
	 ** A pseudo-table created by this opcode is used to hold a single
	 ** row output from the sorter so that the row can be decomposed into
	 ** individual columns using the OP_Column opcode.  The OP_Column opcode
	 ** is the only cursor opcode that works with a pseudo-table.
	 **
	 ** P3 is the number of fields in the records that will be stored by
	 ** the pseudo-table.
	 */
	Op_OpenPseudo = 0x4f

	/* Opcode: Close P1 * * * *
	**
	** Close a cursor previously opened as P1.  If P1 is not
	** currently open, this instruction is a no-op.
	 */
	Op_Close = 0x50

	/* Opcode: ColumnsUsed P1 * * P4 *
	**
	** This opcode (which only exists if SQLite was compiled with
	** SQLITE_ENABLE_COLUMN_USED_MASK) identifies which columns of the
	** table or index for cursor P1 are used.  P4 is a 64-bit integer
	** (P4_INT64) in which the first 63 bits are one for each of the
	** first 63 columns of the table or index that are actually used
	** by the cursor.  The high-order bit is set if any column after
	** the 64th is used.
	 */
	Op_ColumnUsed = 0x51

	/* Opcode: SeekGE P1 P2 P3 P4 *
	 ** Synopsis: key=r[P3@P4]
	 **
	 ** If cursor P1 refers to an SQL table (B-Tree that uses integer keys),
	 ** use the value in register P3 as the key.  If cursor P1 refers
	 ** to an SQL index, then P3 is the first in an array of P4 registers
	 ** that are used as an unpacked index key.
	 **
	 ** Reposition cursor P1 so that  it points to the smallest entry that
	 ** is greater than or equal to the key value. If there are no records
	 ** greater than or equal to the key and P2 is not zero, then jump to P2.
	 **
	 ** If the cursor P1 was opened using the OPFLAG_SEEKEQ flag, then this
	 ** opcode will either land on a record that exactly matches the key, or
	 ** else it will cause a jump to P2.  When the cursor is OPFLAG_SEEKEQ,
	 ** this opcode must be followed by an IdxLE opcode with the same arguments.
	 ** The IdxGT opcode will be skipped if this opcode succeeds, but the
	 ** IdxGT opcode will be used on subsequent loop iterations.  The
	 ** OPFLAG_SEEKEQ flags is a hint to the btree layer to say that this
	 ** is an equality search.
	 **
	 ** This opcode leaves the cursor configured to move in forward order,
	 ** from the beginning toward the end.  In other words, the cursor is
	 ** configured to use Next, not Prev.
	 **
	 ** See also: Found, NotFound, SeekLt, SeekGt, SeekLe
	 */
	Op_SeekGE = 0x52

	/* Opcode: SeekGT P1 P2 P3 P4 *
	** Synopsis: key=r[P3@P4]
	**
	** If cursor P1 refers to an SQL table (B-Tree that uses integer keys),
	** use the value in register P3 as a key. If cursor P1 refers
	** to an SQL index, then P3 is the first in an array of P4 registers
	** that are used as an unpacked index key.
	**
	** Reposition cursor P1 so that it points to the smallest entry that
	** is greater than the key value. If there are no records greater than
	** the key and P2 is not zero, then jump to P2.
	**
	** This opcode leaves the cursor configured to move in forward order,
	** from the beginning toward the end.  In other words, the cursor is
	** configured to use Next, not Prev.
	**
	** See also: Found, NotFound, SeekLt, SeekGe, SeekLe
	 */
	Op_SeekGT = 0x53

	/* Opcode: SeekLT P1 P2 P3 P4 *
	** Synopsis: key=r[P3@P4]
	**
	** If cursor P1 refers to an SQL table (B-Tree that uses integer keys),
	** use the value in register P3 as a key. If cursor P1 refers
	** to an SQL index, then P3 is the first in an array of P4 registers
	** that are used as an unpacked index key.
	**
	** Reposition cursor P1 so that  it points to the largest entry that
	** is less than the key value. If there are no records less than
	** the key and P2 is not zero, then jump to P2.
	**
	** This opcode leaves the cursor configured to move in reverse order,
	** from the end toward the beginning.  In other words, the cursor is
	** configured to use Prev, not Next.
	**
	** See also: Found, NotFound, SeekGt, SeekGe, SeekLe
	 */
	Op_SeekLT = 0x54

	/* Opcode: SeekLE P1 P2 P3 P4 *
	** Synopsis: key=r[P3@P4]
	**
	** If cursor P1 refers to an SQL table (B-Tree that uses integer keys),
	** use the value in register P3 as a key. If cursor P1 refers
	** to an SQL index, then P3 is the first in an array of P4 registers
	** that are used as an unpacked index key.
	**
	** Reposition cursor P1 so that it points to the largest entry that
	** is less than or equal to the key value. If there are no records
	** less than or equal to the key and P2 is not zero, then jump to P2.
	**
	** This opcode leaves the cursor configured to move in reverse order,
	** from the end toward the beginning.  In other words, the cursor is
	** configured to use Prev, not Next.
	**
	** If the cursor P1 was opened using the OPFLAG_SEEKEQ flag, then this
	** opcode will either land on a record that exactly matches the key, or
	** else it will cause a jump to P2.  When the cursor is OPFLAG_SEEKEQ,
	** this opcode must be followed by an IdxLE opcode with the same arguments.
	** The IdxGE opcode will be skipped if this opcode succeeds, but the
	** IdxGE opcode will be used on subsequent loop iterations.  The
	** OPFLAG_SEEKEQ flags is a hint to the btree layer to say that this
	** is an equality search.
	**
	** See also: Found, NotFound, SeekGt, SeekGe, SeekLt
	 */
	Op_SeekLE = 0x55

	/* Opcode: SeekScan  P1 P2 * * *
	** Synopsis: Scan-ahead up to P1 rows
	**
	** This opcode is a prefix opcode to OP_SeekGE.  In other words, this
	** opcode must be immediately followed by OP_SeekGE. This constraint is
	** checked by assert() statements.
	**
	** This opcode uses the P1 through P4 operands of the subsequent
	** OP_SeekGE.  In the text that follows, the operands of the subsequent
	** OP_SeekGE opcode are denoted as SeekOP.P1 through SeekOP.P4.   Only
	** the P1 and P2 operands of this opcode are also used, and  are called
	** This.P1 and This.P2.
	**
	** This opcode helps to optimize IN operators on a multi-column index
	** where the IN operator is on the later terms of the index by avoiding
	** unnecessary seeks on the btree, substituting steps to the next row
	** of the b-tree instead.  A correct answer is obtained if this opcode
	** is omitted or is a no-op.
	**
	** The SeekGE.P3 and SeekGE.P4 operands identify an unpacked key which
	** is the desired entry that we want the cursor SeekGE.P1 to be pointing
	** to.  Call this SeekGE.P4/P5 row the "target".
	**
	** If the SeekGE.P1 cursor is not currently pointing to a valid row,
	** then this opcode is a no-op and control passes through into the OP_SeekGE.
	**
	** If the SeekGE.P1 cursor is pointing to a valid row, then that row
	** might be the target row, or it might be near and slightly before the
	** target row.  This opcode attempts to position the cursor on the target
	** row by, perhaps by invoking sqlite3BtreeStep() on the cursor
	** between 0 and This.P1 times.
	**
	** There are three possible outcomes from this opcode:<ol>
	**
	** <li> If after This.P1 steps, the cursor is still pointing to a place that
	**      is earlier in the btree than the target row, then fall through
	**      into the subsquence OP_SeekGE opcode.
	**
	** <li> If the cursor is successfully moved to the target row by 0 or more
	**      sqlite3BtreeNext() calls, then jump to This.P2, which will land just
	**      past the OP_IdxGT or OP_IdxGE opcode that follows the OP_SeekGE.
	**
	** <li> If the cursor ends up past the target row (indicating the the target
	**      row does not exist in the btree) then jump to SeekOP.P2.
	** </ol>
	 */
	Op_SeekScan = 0x56

	/* Opcode: SeekHit P1 P2 P3 * *
	 ** Synopsis: set P2<=seekHit<=P3
	 **
	 ** Increase or decrease the seekHit value for cursor P1, if necessary,
	 ** so that it is no less than P2 and no greater than P3.
	 **
	 ** The seekHit integer represents the maximum of terms in an index for which
	 ** there is known to be at least one match.  If the seekHit value is smaller
	 ** than the total number of equality terms in an index lookup, then the
	 ** OP_IfNoHope opcode might run to see if the IN loop can be abandoned
	 ** early, thus saving work.  This is part of the IN-early-out optimization.
	 **
	 ** P1 must be a valid b-tree cursor.
	 */
	Op_SeekHit = 0x57

	/* Opcode: IfNotOpen P1 P2 * * *
	** Synopsis: if( !csr[P1] ) goto P2
	**
	** If cursor P1 is not open, jump to instruction P2. Otherwise, fall through.
	 */
	Op_IfNotOpen = 0x58

	/* Opcode: Found P1 P2 P3 P4 *
	** Synopsis: key=r[P3@P4]
	**
	** If P4==0 then register P3 holds a blob constructed by MakeRecord.  If
	** P4>0 then register P3 is the first of P4 registers that form an unpacked
	** record.
	**
	** Cursor P1 is on an index btree.  If the record identified by P3 and P4
	** is a prefix of any entry in P1 then a jump is made to P2 and
	** P1 is left pointing at the matching entry.
	**
	** This operation leaves the cursor in a state where it can be
	** advanced in the forward direction.  The Next instruction will work,
	** but not the Prev instruction.
	**
	** See also: NotFound, NoConflict, NotExists. SeekGe
	 */
	Op_Found = 0x59

	/* Opcode: NotFound P1 P2 P3 P4 *
	** Synopsis: key=r[P3@P4]
	**
	** If P4==0 then register P3 holds a blob constructed by MakeRecord.  If
	** P4>0 then register P3 is the first of P4 registers that form an unpacked
	** record.
	**
	** Cursor P1 is on an index btree.  If the record identified by P3 and P4
	** is not the prefix of any entry in P1 then a jump is made to P2.  If P1
	** does contain an entry whose prefix matches the P3/P4 record then control
	** falls through to the next instruction and P1 is left pointing at the
	** matching entry.
	**
	** This operation leaves the cursor in a state where it cannot be
	** advanced in either direction.  In other words, the Next and Prev
	** opcodes do not work after this operation.
	**
	** See also: Found, NotExists, NoConflict, IfNoHope
	 */
	Op_NotFound = 0x5a

	/* Opcode: IfNoHope P1 P2 P3 P4 *
	** Synopsis: key=r[P3@P4]
	**
	** Register P3 is the first of P4 registers that form an unpacked
	** record.  Cursor P1 is an index btree.  P2 is a jump destination.
	** In other words, the operands to this opcode are the same as the
	** operands to OP_NotFound and OP_IdxGT.
	**
	** This opcode is an optimization attempt only.  If this opcode always
	** falls through, the correct answer is still obtained, but extra works
	** is performed.
	**
	** A value of N in the seekHit flag of cursor P1 means that there exists
	** a key P3:N that will match some record in the index.  We want to know
	** if it is possible for a record P3:P4 to match some record in the
	** index.  If it is not possible, we can skips some work.  So if seekHit
	** is less than P4, attempt to find out if a match is possible by running
	** OP_NotFound.
	**
	** This opcode is used in IN clause processing for a multi-column key.
	** If an IN clause is attached to an element of the key other than the
	** left-most element, and if there are no matches on the most recent
	** seek over the whole key, then it might be that one of the key element
	** to the left is prohibiting a match, and hence there is "no hope" of
	** any match regardless of how many IN clause elements are checked.
	** In such a case, we abandon the IN clause search early, using this
	** opcode.  The opcode name comes from the fact that the
	** jump is taken if there is "no hope" of achieving a match.
	**
	** See also: NotFound, SeekHit
	 */
	Op_IfNoHope = 0x5b

	/* Opcode: NoConflict P1 P2 P3 P4 *
	** Synopsis: key=r[P3@P4]
	**
	** If P4==0 then register P3 holds a blob constructed by MakeRecord.  If
	** P4>0 then register P3 is the first of P4 registers that form an unpacked
	** record.
	**
	** Cursor P1 is on an index btree.  If the record identified by P3 and P4
	** contains any NULL value, jump immediately to P2.  If all terms of the
	** record are not-NULL then a check is done to determine if any row in the
	** P1 index btree has a matching key prefix.  If there are no matches, jump
	** immediately to P2.  If there is a match, fall through and leave the P1
	** cursor pointing to the matching row.
	**
	** This opcode is similar to OP_NotFound with the exceptions that the
	** branch is always taken if any part of the search key input is NULL.
	**
	** This operation leaves the cursor in a state where it cannot be
	** advanced in either direction.  In other words, the Next and Prev
	** opcodes do not work after this operation.
	**
	** See also: NotFound, Found, NotExists
	 */
	Op_NoConflict = 0x5c

	/* Opcode: SeekRowid P1 P2 P3 * *
	** Synopsis: intkey=r[P3]
	**
	** P1 is the index of a cursor open on an SQL table btree (with integer
	** keys).  If register P3 does not contain an integer or if P1 does not
	** contain a record with rowid P3 then jump immediately to P2.
	** Or, if P2 is 0, raise an SQLITE_CORRUPT error. If P1 does contain
	** a record with rowid P3 then
	** leave the cursor pointing at that record and fall through to the next
	** instruction.
	**
	** The OP_NotExists opcode performs the same operation, but with OP_NotExists
	** the P3 register must be guaranteed to contain an integer value.  With this
	** opcode, register P3 might not contain an integer.
	**
	** The OP_NotFound opcode performs the same operation on index btrees
	** (with arbitrary multi-value keys).
	**
	** This opcode leaves the cursor in a state where it cannot be advanced
	** in either direction.  In other words, the Next and Prev opcodes will
	** not work following this opcode.
	**
	** See also: Found, NotFound, NoConflict, SeekRowid
	 */
	Op_SeekRowid = 0x5d

	/* Opcode: NotExists P1 P2 P3 * *
	** Synopsis: intkey=r[P3]
	**
	** P1 is the index of a cursor open on an SQL table btree (with integer
	** keys).  P3 is an integer rowid.  If P1 does not contain a record with
	** rowid P3 then jump immediately to P2.  Or, if P2 is 0, raise an
	** SQLITE_CORRUPT error. If P1 does contain a record with rowid P3 then
	** leave the cursor pointing at that record and fall through to the next
	** instruction.
	**
	** The OP_SeekRowid opcode performs the same operation but also allows the
	** P3 register to contain a non-integer value, in which case the jump is
	** always taken.  This opcode requires that P3 always contain an integer.
	**
	** The OP_NotFound opcode performs the same operation on index btrees
	** (with arbitrary multi-value keys).
	**
	** This opcode leaves the cursor in a state where it cannot be advanced
	** in either direction.  In other words, the Next and Prev opcodes will
	** not work following this opcode.
	**
	** See also: Found, NotFound, NoConflict, SeekRowid
	 */
	OP_NotExists = 0x5e

	/* Opcode: Sequence P1 P2 * * *
	** Synopsis: r[P2]=cursor[P1].ctr++
	**
	** Find the next available sequence number for cursor P1.
	** Write the sequence number into register P2.
	** The sequence number on the cursor is incremented after this
	** instruction.
	 */
	Op_Sequence = 0x5f

	/* Opcode: NewRowid P1 P2 P3 * *
	** Synopsis: r[P2]=rowid
	**
	** Get a new integer record number (a.k.a "rowid") used as the key to a table.
	** The record number is not previously used as a key in the database
	** table that cursor P1 points to.  The new record number is written
	** written to register P2.
	**
	** If P3>0 then P3 is a register in the root frame of this VDBE that holds
	** the largest previously generated record number. No new record numbers are
	** allowed to be less than this value. When this value reaches its maximum,
	** an SQLITE_FULL error is generated. The P3 register is updated with the '
	** generated record number. This P3 mechanism is used to help implement the
	** AUTOINCREMENT feature.
	 */
	Op_NewRowid = 0x60

	/* Opcode: Insert P1 P2 P3 P4 P5
	** Synopsis: intkey=r[P3] data=r[P2]
	**
	** Write an entry into the table of cursor P1.  A new entry is
	** created if it doesn't already exist or the data for an existing
	** entry is overwritten.  The data is the value MEM_Blob stored in register
	** number P2. The key is stored in register P3. The key must
	** be a MEM_Int.
	**
	** If the OPFLAG_NCHANGE flag of P5 is set, then the row change count is
	** incremented (otherwise not).  If the OPFLAG_LASTROWID flag of P5 is set,
	** then rowid is stored for subsequent return by the
	** sqlite3_last_insert_rowid() function (otherwise it is unmodified).
	**
	** If the OPFLAG_USESEEKRESULT flag of P5 is set, the implementation might
	** run faster by avoiding an unnecessary seek on cursor P1.  However,
	** the OPFLAG_USESEEKRESULT flag must only be set if there have been no prior
	** seeks on the cursor or if the most recent seek used a key equal to P3.
	**
	** If the OPFLAG_ISUPDATE flag is set, then this opcode is part of an
	** UPDATE operation.  Otherwise (if the flag is clear) then this opcode
	** is part of an INSERT operation.  The difference is only important to
	** the update hook.
	**
	** Parameter P4 may point to a Table structure, or may be NULL. If it is
	** not NULL, then the update-hook (sqlite3.xUpdateCallback) is invoked
	** following a successful insert.
	**
	** (WARNING/TODO: If P1 is a pseudo-cursor and P2 is dynamically
	** allocated, then ownership of P2 is transferred to the pseudo-cursor
	** and register P2 becomes ephemeral.  If the cursor is changed, the
	** value of register P2 will then change.  Make sure this does not
	** cause any problems.)
	**
	** This instruction only works on tables.  The equivalent instruction
	** for indices is OP_IdxInsert.
	 */
	Op_Insert = 0x61

	/* Opcode: RowCell P1 P2 P3 * *
	 **
	 ** P1 and P2 are both open cursors. Both must be opened on the same type
	 ** of table - intkey or index. This opcode is used as part of copying
	 ** the current row from P2 into P1. If the cursors are opened on intkey
	 ** tables, register P3 contains the rowid to use with the new record in
	 ** P1. If they are opened on index tables, P3 is not used.
	 **
	 ** This opcode must be followed by either an Insert or InsertIdx opcode
	 ** with the OPFLAG_PREFORMAT flag set to complete the insert operation.
	 */
	Op_RowCell = 0x62

	/* Opcode: Delete P1 P2 P3 P4 P5
	 **
	 ** Delete the record at which the P1 cursor is currently pointing.
	 **
	 ** If the OPFLAG_SAVEPOSITION bit of the P5 parameter is set, then
	 ** the cursor will be left pointing at  either the next or the previous
	 ** record in the table. If it is left pointing at the next record, then
	 ** the next Next instruction will be a no-op. As a result, in this case
	 ** it is ok to delete a record from within a Next loop. If
	 ** OPFLAG_SAVEPOSITION bit of P5 is clear, then the cursor will be
	 ** left in an undefined state.
	 **
	 ** If the OPFLAG_AUXDELETE bit is set on P5, that indicates that this
	 ** delete one of several associated with deleting a table row and all its
	 ** associated index entries.  Exactly one of those deletes is the "primary"
	 ** delete.  The others are all on OPFLAG_FORDELETE cursors or else are
	 ** marked with the AUXDELETE flag.
	 **
	 ** If the OPFLAG_NCHANGE flag of P2 (NB: P2 not P5) is set, then the row
	 ** change count is incremented (otherwise not).
	 **
	 ** P1 must not be pseudo-table.  It has to be a real table with
	 ** multiple rows.
	 **
	 ** If P4 is not NULL then it points to a Table object. In this case either
	 ** the update or pre-update hook, or both, may be invoked. The P1 cursor must
	 ** have been positioned using OP_NotFound prior to invoking this opcode in
	 ** this case. Specifically, if one is configured, the pre-update hook is
	 ** invoked if P4 is not NULL. The update-hook is invoked if one is configured,
	 ** P4 is not NULL, and the OPFLAG_NCHANGE flag is set in P2.
	 **
	 ** If the OPFLAG_ISUPDATE flag is set in P2, then P3 contains the address
	 ** of the memory cell that contains the value that the rowid of the row will
	 ** be set to by the update.
	 */
	Op_Delete = 0x63

	/* Opcode: ResetCount * * * * *
	**
	** The value of the change counter is copied to the database handle
	** change counter (returned by subsequent calls to sqlite3_changes()).
	** Then the VMs internal change counter resets to 0.
	** This is used by trigger programs.
	 */
	Op_ResetCount = 0x64

	/* Opcode: SorterCompare P1 P2 P3 P4
	** Synopsis: if key(P1)!=trim(r[P3],P4) goto P2
	**
	** P1 is a sorter cursor. This instruction compares a prefix of the
	** record blob in register P3 against a prefix of the entry that
	** the sorter cursor currently points to.  Only the first P4 fields
	** of r[P3] and the sorter record are compared.
	**
	** If either P3 or the sorter contains a NULL in one of their significant
	** fields (not counting the P4 fields at the end which are ignored) then
	** the comparison is assumed to be equal.
	**
	** Fall through to next instruction if the two records compare equal to
	** each other.  Jump to P2 if they are different.
	 */
	Op_SorterCpmpare = 0x65

	/* Opcode: SorterData P1 P2 P3 * *
	 ** Synopsis: r[P2]=data
	 **
	 ** Write into register P2 the current sorter data for sorter cursor P1.
	 ** Then clear the column header cache on cursor P3.
	 **
	 ** This opcode is normally use to move a record out of the sorter and into
	 ** a register that is the source for a pseudo-table cursor created using
	 ** OpenPseudo.  That pseudo-table cursor is the one that is identified by
	 ** parameter P3.  Clearing the P3 column cache as part of this opcode saves
	 ** us from having to issue a separate NullRow instruction to clear that cache.
	 */
	Op_SorterData = 0x66

	/* Opcode: RowData P1 P2 P3 * *
	** Synopsis: r[P2]=data
	**
	** Write into register P2 the complete row content for the row at
	** which cursor P1 is currently pointing.
	** There is no interpretation of the data.
	** It is just copied onto the P2 register exactly as
	** it is found in the database file.
	**
	** If cursor P1 is an index, then the content is the key of the row.
	** If cursor P2 is a table, then the content extracted is the data.
	**
	** If the P1 cursor must be pointing to a valid row (not a NULL row)
	** of a real table, not a pseudo-table.
	**
	** If P3!=0 then this opcode is allowed to make an ephemeral pointer
	** into the database page.  That means that the content of the output
	** register will be invalidated as soon as the cursor moves - including
	** moves caused by other cursors that "save" the current cursors
	** position in order that they can write to the same table.  If P3==0
	** then a copy of the data is made into memory.  P3!=0 is faster, but
	** P3==0 is safer.
	**
	** If P3!=0 then the content of the P2 register is unsuitable for use
	** in OP_Result and any OP_Result will invalidate the P2 register content.
	** The P2 register content is invalidated by opcodes like OP_Function or
	** by any use of another cursor pointing to the same table.
	 */
	Op_RowData = 0x67

	/* Opcode: Rowid P1 P2 * * *
	** Synopsis: r[P2]=rowid
	**
	** Store in register P2 an integer which is the key of the table entry that
	** P1 is currently point to.
	**
	** P1 can be either an ordinary table or a virtual table.  There used to
	** be a separate OP_VRowid opcode for use with virtual tables, but this
	** one opcode now works for both table types.
	 */
	Op_Rowid = 0x68

	/* Opcode: NullRow P1 * * * *
	**
	** Move the cursor P1 to a null row.  Any OP_Column operations
	** that occur while the cursor is on the null row will always
	** write a NULL.
	 */
	Op_NullRow = 0x69

	/* Opcode: SeekEnd P1 * * * *
	**
	** Position cursor P1 at the end of the btree for the purpose of
	** appending a new entry onto the btree.
	**
	** It is assumed that the cursor is used only for appending and so
	** if the cursor is valid, then the cursor must already be pointing
	** at the end of the btree and so no changes are made to
	** the cursor.
	 */
	Op_SeekEnd = 0x6a

	/* Opcode: Last P1 P2 * * *
	**
	** The next use of the Rowid or Column or Prev instruction for P1
	** will refer to the last entry in the database table or index.
	** If the table or index is empty and P2>0, then jump immediately to P2.
	** If P2 is 0 or if the table or index is not empty, fall through
	** to the following instruction.
	**
	** This opcode leaves the cursor configured to move in reverse order,
	** from the end toward the beginning.  In other words, the cursor is
	** configured to use Prev, not Next.
	 */
	Op_Last = 0x6b

	/* Opcode: IfSmaller P1 P2 P3 * *
	**
	** Estimate the number of rows in the table P1.  Jump to P2 if that
	** estimate is less than approximately 2**(0.1*P3).
	 */
	Op_IfSmaller = 0x6c

	/* Opcode: SorterSort P1 P2 * * *
	**
	** After all records have been inserted into the Sorter object
	** identified by P1, invoke this opcode to actually do the sorting.
	** Jump to P2 if there are no records to be sorted.
	**
	** This opcode is an alias for OP_Sort and OP_Rewind that is used
	** for Sorter objects.
	 */
	Op_SorterSort = 0x6d

	/* Opcode: Sort P1 P2 * * *
	**
	** This opcode does exactly the same thing as OP_Rewind except that
	** it increments an undocumented global variable used for testing.
	**
	** Sorting is accomplished by writing records into a sorting index,
	** then rewinding that index and playing it back from beginning to
	** end.  We use the OP_Sort opcode instead of OP_Rewind to do the
	** rewinding so that the global variable will be incremented and
	** regression tests can determine whether or not the optimizer is
	** correctly optimizing out sorts.
	 */
	Op_Sort = 0x6e

	/* Opcode: Rewind P1 P2 * * *
	**
	** The next use of the Rowid or Column or Next instruction for P1
	** will refer to the first entry in the database table or index.
	** If the table or index is empty, jump immediately to P2.
	** If the table or index is not empty, fall through to the following
	** instruction.
	**
	** This opcode leaves the cursor configured to move in forward order,
	** from the beginning toward the end.  In other words, the cursor is
	** configured to use Next, not Prev.
	 */
	Op_Rewind = 0x6f

	/* Opcode: Next P1 P2 P3 P4 P5
	**
	** Advance cursor P1 so that it points to the next key/data pair in its
	** table or index.  If there are no more key/value pairs then fall through
	** to the following instruction.  But if the cursor advance was successful,
	** jump immediately to P2.
	**
	** The Next opcode is only valid following an SeekGT, SeekGE, or
	** OP_Rewind opcode used to position the cursor.  Next is not allowed
	** to follow SeekLT, SeekLE, or OP_Last.
	**
	** The P1 cursor must be for a real table, not a pseudo-table.  P1 must have
	** been opened prior to this opcode or the program will segfault.
	**
	** The P3 value is a hint to the btree implementation. If P3==1, that
	** means P1 is an SQL index and that this instruction could have been
	** omitted if that index had been unique.  P3 is usually 0.  P3 is
	** always either 0 or 1.
	**
	** P4 is always of type P4_ADVANCE. The function pointer points to
	** sqlite3BtreeNext().
	**
	** If P5 is positive and the jump is taken, then event counter
	** number P5-1 in the prepared statement is incremented.
	**
	** See also: Prev
	 */
	Op_Next = 0x70

	/* Opcode: Prev P1 P2 P3 P4 P5
	**
	** Back up cursor P1 so that it points to the previous key/data pair in its
	** table or index.  If there is no previous key/value pairs then fall through
	** to the following instruction.  But if the cursor backup was successful,
	** jump immediately to P2.
	**
	**
	** The Prev opcode is only valid following an SeekLT, SeekLE, or
	** OP_Last opcode used to position the cursor.  Prev is not allowed
	** to follow SeekGT, SeekGE, or OP_Rewind.
	**
	** The P1 cursor must be for a real table, not a pseudo-table.  If P1 is
	** not open then the behavior is undefined.
	**
	** The P3 value is a hint to the btree implementation. If P3==1, that
	** means P1 is an SQL index and that this instruction could have been
	** omitted if that index had been unique.  P3 is usually 0.  P3 is
	** always either 0 or 1.
	**
	** P4 is always of type P4_ADVANCE. The function pointer points to
	** sqlite3BtreePrevious().
	**
	** If P5 is positive and the jump is taken, then event counter
	** number P5-1 in the prepared statement is incremented.
	 */
	Op_Prev = 0x71

	/* Opcode: SorterNext P1 P2 * * P5
	**
	** This opcode works just like OP_Next except that P1 must be a
	** sorter object for which the OP_SorterSort opcode has been
	** invoked.  This opcode advances the cursor to the next sorted
	** record, or jumps to P2 if there are no more sorted records.
	 */
	Op_SorterNext = 0x72

	/* Opcode: IdxInsert P1 P2 P3 P4 P5
	 ** Synopsis: key=r[P2]
	 **
	 ** Register P2 holds an SQL index key made using the
	 ** MakeRecord instructions.  This opcode writes that key
	 ** into the index P1.  Data for the entry is nil.
	 **
	 ** If P4 is not zero, then it is the number of values in the unpacked
	 ** key of reg(P2).  In that case, P3 is the index of the first register
	 ** for the unpacked key.  The availability of the unpacked key can sometimes
	 ** be an optimization.
	 **
	 ** If P5 has the OPFLAG_APPEND bit set, that is a hint to the b-tree layer
	 ** that this insert is likely to be an append.
	 **
	 ** If P5 has the OPFLAG_NCHANGE bit set, then the change counter is
	 ** incremented by this instruction.  If the OPFLAG_NCHANGE bit is clear,
	 ** then the change counter is unchanged.
	 **
	 ** If the OPFLAG_USESEEKRESULT flag of P5 is set, the implementation might
	 ** run faster by avoiding an unnecessary seek on cursor P1.  However,
	 ** the OPFLAG_USESEEKRESULT flag must only be set if there have been no prior
	 ** seeks on the cursor or if the most recent seek used a key equivalent
	 ** to P2.
	 **
	 ** This instruction only works for indices.  The equivalent instruction
	 ** for tables is OP_Insert.
	 */
	Op_IdxInsert = 0x73

	/* Opcode: SorterInsert P1 P2 * * *
	 ** Synopsis: key=r[P2]
	 **
	 ** Register P2 holds an SQL index key made using the
	 ** MakeRecord instructions.  This opcode writes that key
	 ** into the sorter P1.  Data for the entry is nil.
	 */
	Op_SorterInsert = 0x74

	/* Opcode: IdxDelete P1 P2 P3 * P5
	** Synopsis: key=r[P2@P3]
	**
	** The content of P3 registers starting at register P2 form
	** an unpacked index key. This opcode removes that entry from the
	** index opened by cursor P1.
	**
	** If P5 is not zero, then raise an SQLITE_CORRUPT_INDEX error
	** if no matching index entry is found.  This happens when running
	** an UPDATE or DELETE statement and the index entry to be updated
	** or deleted is not found.  For some uses of IdxDelete
	** (example:  the EXCEPT operator) it does not matter that no matching
	** entry is found.  For those cases, P5 is zero.
	 */
	Op_IdxDelete = 0x75

	/* Opcode: DeferredSeek P1 * P3 P4 *
	** Synopsis: Move P3 to P1.rowid if needed
	**
	** P1 is an open index cursor and P3 is a cursor on the corresponding
	** table.  This opcode does a deferred seek of the P3 table cursor
	** to the row that corresponds to the current row of P1.
	**
	** This is a deferred seek.  Nothing actually happens until
	** the cursor is used to read a record.  That way, if no reads
	** occur, no unnecessary I/O happens.
	**
	** P4 may be an array of integers (type P4_INTARRAY) containing
	** one entry for each column in the P3 table.  If array entry a(i)
	** is non-zero, then reading column a(i)-1 from cursor P3 is
	** equivalent to performing the deferred seek and then reading column i
	** from P1.  This information is stored in P3 and used to redirect
	** reads against P3 over to P1, thus possibly avoiding the need to
	** seek and read cursor P3.
	 */
	Op_DeferredSeek = 0x76

	/* Opcode: IdxRowid P1 P2 * * *
	** Synopsis: r[P2]=rowid
	**
	** Write into register P2 an integer which is the last entry in the record at
	** the end of the index key pointed to by cursor P1.  This integer should be
	** the rowid of the table entry to which this index entry points.
	**
	** See also: Rowid, MakeRecord.
	 */
	Op_IdxRowid = 0x77

	/* Opcode: FinishSeek P1 * * * *
	 **
	 ** If cursor P1 was previously moved via OP_DeferredSeek, complete that
	 ** seek operation now, without further delay.  If the cursor seek has
	 ** already occurred, this instruction is a no-op.
	 */
	Op_FinishSeek = 0x78

	/* Opcode: IdxGE P1 P2 P3 P4 *
	** Synopsis: key=r[P3@P4]
	**
	** The P4 register values beginning with P3 form an unpacked index
	** key that omits the PRIMARY KEY.  Compare this key value against the index
	** that P1 is currently pointing to, ignoring the PRIMARY KEY or ROWID
	** fields at the end.
	**
	** If the P1 index entry is greater than or equal to the key value
	** then jump to P2.  Otherwise fall through to the next instruction.
	 */
	Op_IdxGE = 0x79

	/* Opcode: IdxGT P1 P2 P3 P4 *
	** Synopsis: key=r[P3@P4]
	**
	** The P4 register values beginning with P3 form an unpacked index
	** key that omits the PRIMARY KEY.  Compare this key value against the index
	** that P1 is currently pointing to, ignoring the PRIMARY KEY or ROWID
	** fields at the end.
	**
	** If the P1 index entry is greater than the key value
	** then jump to P2.  Otherwise fall through to the next instruction.
	 */
	Op_IdxGT = 0x7a

	/* Opcode: IdxLT P1 P2 P3 P4 *
	** Synopsis: key=r[P3@P4]
	**
	** The P4 register values beginning with P3 form an unpacked index
	** key that omits the PRIMARY KEY or ROWID.  Compare this key value against
	** the index that P1 is currently pointing to, ignoring the PRIMARY KEY or
	** ROWID on the P1 index.
	**
	** If the P1 index entry is less than the key value then jump to P2.
	** Otherwise fall through to the next instruction.
	 */
	Op_IdxLT = 0x7b

	/* Opcode: IdxLE P1 P2 P3 P4 *
	** Synopsis: key=r[P3@P4]
	**
	** The P4 register values beginning with P3 form an unpacked index
	** key that omits the PRIMARY KEY or ROWID.  Compare this key value against
	** the index that P1 is currently pointing to, ignoring the PRIMARY KEY or
	** ROWID on the P1 index.
	**
	** If the P1 index entry is less than or equal to the key value then jump
	** to P2. Otherwise fall through to the next instruction.
	 */
	Op_IdxLE = 0x7c

	/* Opcode: Destroy P1 P2 P3 * *
	**
	** Delete an entire database table or index whose root page in the database
	** file is given by P1.
	**
	** The table being destroyed is in the main database file if P3==0.  If
	** P3==1 then the table to be clear is in the auxiliary database file
	** that is used to store tables create using CREATE TEMPORARY TABLE.
	**
	** If AUTOVACUUM is enabled then it is possible that another root page
	** might be moved into the newly deleted root page in order to keep all
	** root pages contiguous at the beginning of the database.  The former
	** value of the root page that moved - its value before the move occurred -
	** is stored in register P2. If no page movement was required (because the
	** table being dropped was already the last one in the database) then a
	** zero is stored in register P2.  If AUTOVACUUM is disabled then a zero
	** is stored in register P2.
	**
	** This opcode throws an error if there are any active reader VMs when
	** it is invoked. This is done to avoid the difficulty associated with
	** updating existing cursors when a root page is moved in an AUTOVACUUM
	** database. This error is thrown even if the database is not an AUTOVACUUM
	** db in order to avoid introducing an incompatibility between autovacuum
	** and non-autovacuum modes.
	**
	** See also: Clear
	 */
	Op_Destroy = 0x7d

	/* Opcode: Clear P1 P2 P3
	**
	** Delete all contents of the database table or index whose root page
	** in the database file is given by P1.  But, unlike Destroy, do not
	** remove the table or index from the database file.
	**
	** The table being clear is in the main database file if P2==0.  If
	** P2==1 then the table to be clear is in the auxiliary database file
	** that is used to store tables create using CREATE TEMPORARY TABLE.
	**
	** If the P3 value is non-zero, then the row change count is incremented
	** by the number of rows in the table being cleared. If P3 is greater
	** than zero, then the value stored in register P3 is also incremented
	** by the number of rows in the table being cleared.
	**
	** See also: Destroy
	 */
	Op_Clear = 0x7e

	/* Opcode: ResetSorter P1 * * * *
	 **
	 ** Delete all contents from the ephemeral table or sorter
	 ** that is open on cursor P1.
	 **
	 ** This opcode only works for cursors used for sorting and
	 ** opened with OP_OpenEphemeral or OP_SorterOpen.
	 */
	Op_ResetSorter = 0x7f

	/* Opcode: CreateBtree P1 P2 P3 * *
	** Synopsis: r[P2]=root iDb=P1 flags=P3
	**
	** Allocate a new b-tree in the main database file if P1==0 or in the
	** TEMP database file if P1==1 or in an attached database if
	** P1>1.  The P3 argument must be 1 (BTREE_INTKEY) for a rowid table
	** it must be 2 (BTREE_BLOBKEY) for an index or WITHOUT ROWID table.
	** The root page number of the new b-tree is stored in register P2.
	 */
	Op_CreateBtree = 0x80

	/* Opcode: SqlExec * * * P4 *
	**
	** Run the SQL statement or statements specified in the P4 string.
	 */
	Op_SqlExec = 0x81

	/* Opcode: ParseSchema P1 * * P4 *
	**
	** Read and parse all entries from the schema table of database P1
	** that match the WHERE clause P4.  If P4 is a NULL pointer, then the
	** entire schema for P1 is reparsed.
	**
	** This opcode invokes the parser to create a new virtual machine,
	** then runs the new virtual machine.  It is thus a re-entrant opcode.
	 */
	Op_ParseSchema = 0x82

	/* Opcode: LoadAnalysis P1 * * * *
	**
	** Read the sqlite_stat1 table for database P1 and load the content
	** of that table into the internal index hash table.  This will cause
	** the analysis to be used when preparing all subsequent queries.
	 */
	Op_LoadAnalysis = 0x83

	/* Opcode: DropTable P1 * * P4 *
	**
	** Remove the internal (in-memory) data structures that describe
	** the table named P4 in database P1.  This is called after a table
	** is dropped from disk (using the Destroy opcode) in order to keep
	** the internal representation of the
	** schema consistent with what is on disk.
	 */
	Op_DropTable = 0x84

	/* Opcode: DropIndex P1 * * P4 *
	**
	** Remove the internal (in-memory) data structures that describe
	** the index named P4 in database P1.  This is called after an index
	** is dropped from disk (using the Destroy opcode)
	** in order to keep the internal representation of the
	** schema consistent with what is on disk.
	 */
	Op_DropIndex = 0x85

	/* Opcode: DropTrigger P1 * * P4 *
	 **
	 ** Remove the internal (in-memory) data structures that describe
	 ** the trigger named P4 in database P1.  This is called after a trigger
	 ** is dropped from disk (using the Destroy opcode) in order to keep
	 ** the internal representation of the
	 ** schema consistent with what is on disk.
	 */
	Op_DropTrigger = 0x86

	/* Opcode: IntegrityCk P1 P2 P3 P4 P5
	**
	** Do an analysis of the currently open database.  Store in
	** register P1 the text of an error message describing any problems.
	** If no problems are found, store a NULL in register P1.
	**
	** The register P3 contains one less than the maximum number of allowed errors.
	** At most reg(P3) errors will be reported.
	** In other words, the analysis stops as soon as reg(P1) errors are
	** seen.  Reg(P1) is updated with the number of errors remaining.
	**
	** The root page numbers of all tables in the database are integers
	** stored in P4_INTARRAY argument.
	**
	** If P5 is not zero, the check is done on the auxiliary database
	** file, not the main database file.
	**
	** This opcode is used to implement the integrity_check pragma.
	 */
	Op_IntegrityCk = 0x87

	/* Opcode: RowSetAdd P1 P2 * * *
	** Synopsis: rowset(P1)=r[P2]
	**
	** Insert the integer value held by register P2 into a RowSet object
	** held in register P1.
	**
	** An assertion fails if P2 is not an integer.
	 */
	Op_RowSetAdd = 0x88

	/* Opcode: RowSetRead P1 P2 P3 * *
	** Synopsis: r[P3]=rowset(P1)
	**
	** Extract the smallest value from the RowSet object in P1
	** and put that value into register P3.
	** Or, if RowSet object P1 is initially empty, leave P3
	** unchanged and jump to instruction P2.
	 */
	Op_RowSetRead = 0x89

	/* Opcode: RowSetTest P1 P2 P3 P4
	** Synopsis: if r[P3] in rowset(P1) goto P2
	**
	** Register P3 is assumed to hold a 64-bit integer value. If register P1
	** contains a RowSet object and that RowSet object contains
	** the value held in P3, jump to register P2. Otherwise, insert the
	** integer in P3 into the RowSet and continue on to the
	** next opcode.
	**
	** The RowSet object is optimized for the case where sets of integers
	** are inserted in distinct phases, which each set contains no duplicates.
	** Each set is identified by a unique P4 value. The first set
	** must have P4==0, the final set must have P4==-1, and for all other sets
	** must have P4>0.
	**
	** This allows optimizations: (a) when P4==0 there is no need to test
	** the RowSet object for P3, as it is guaranteed not to contain it,
	** (b) when P4==-1 there is no need to insert the value, as it will
	** never be tested for, and (c) when a value that is part of set X is
	** inserted, there is no need to search to see if the same value was
	** previously inserted as part of set X (only if it was previously
	** inserted as part of some other set).
	 */
	Op_RowSetTest = 0x8a

	/* Opcode: Program P1 P2 P3 P4 P5
	 **
	 ** Execute the trigger program passed as P4 (type P4_SUBPROGRAM).
	 **
	 ** P1 contains the address of the memory cell that contains the first memory
	 ** cell in an array of values used as arguments to the sub-program. P2
	 ** contains the address to jump to if the sub-program throws an IGNORE
	 ** exception using the RAISE() function. Register P3 contains the address
	 ** of a memory cell in this (the parent) VM that is used to allocate the
	 ** memory required by the sub-vdbe at runtime.
	 **
	 ** P4 is a pointer to the VM containing the trigger program.
	 **
	 ** If P5 is non-zero, then recursive program invocation is enabled.
	 */
	Op_Program = 0x8b

	/* Opcode: Param P1 P2 * * *
	**
	** This opcode is only ever present in sub-programs called via the
	** OP_Program instruction. Copy a value currently stored in a memory
	** cell of the calling (parent) frame to cell P2 in the current frames
	** address space. This is used by trigger programs to access the new.*
	** and old.* values.
	**
	** The address of the cell in the parent frame is determined by adding
	** the value of the P1 argument to the value of the P1 argument to the
	** calling OP_Program instruction.
	 */
	Op_Param = 0x8c

	/* Opcode: FkCounter P1 P2 * * *
	** Synopsis: fkctr[P1]+=P2
	**
	** Increment a "constraint counter" by P2 (P2 may be negative or positive).
	** If P1 is non-zero, the database constraint counter is incremented
	** (deferred foreign key constraints). Otherwise, if P1 is zero, the
	** statement counter is incremented (immediate foreign key constraints).
	 */
	Op_FkCounter = 0x8d

	/* Opcode: FkIfZero P1 P2 * * *
	** Synopsis: if fkctr[P1]==0 goto P2
	**
	** This opcode tests if a foreign key constraint-counter is currently zero.
	** If so, jump to instruction P2. Otherwise, fall through to the next
	** instruction.
	**
	** If P1 is non-zero, then the jump is taken if the database constraint-counter
	** is zero (the one that counts deferred constraint violations). If P1 is
	** zero, the jump is taken if the statement constraint-counter is zero
	** (immediate foreign key constraint violations).
	 */
	Op_FkIfZero = 0x8e

	/* Opcode: MemMax P1 P2 * * *
	** Synopsis: r[P1]=max(r[P1],r[P2])
	**
	** P1 is a register in the root frame of this VM (the root frame is
	** different from the current frame if this instruction is being executed
	** within a sub-program). Set the value of register P1 to the maximum of
	** its current value and the value in register P2.
	**
	** This instruction throws an error if the memory cell is not initially
	** an integer.
	 */
	Op_MemMax = 0x8f

	/* Opcode: IfPos P1 P2 P3 * *
	** Synopsis: if r[P1]>0 then r[P1]-=P3, goto P2
	**
	** Register P1 must contain an integer.
	** If the value of register P1 is 1 or greater, subtract P3 from the
	** value in P1 and jump to P2.
	**
	** If the initial value of register P1 is less than 1, then the
	** value is unchanged and control passes through to the next instruction.
	 */
	Op_IfPos = 0x90

	/* Opcode: OffsetLimit P1 P2 P3 * *
	** Synopsis: if r[P1]>0 then r[P2]=r[P1]+max(0,r[P3]) else r[P2]=(-1)
	**
	** This opcode performs a commonly used computation associated with
	** LIMIT and OFFSET process.  r[P1] holds the limit counter.  r[P3]
	** holds the offset counter.  The opcode computes the combined value
	** of the LIMIT and OFFSET and stores that value in r[P2].  The r[P2]
	** value computed is the total number of rows that will need to be
	** visited in order to complete the query.
	**
	** If r[P3] is zero or negative, that means there is no OFFSET
	** and r[P2] is set to be the value of the LIMIT, r[P1].
	**
	** if r[P1] is zero or negative, that means there is no LIMIT
	** and r[P2] is set to -1.
	**
	** Otherwise, r[P2] is set to the sum of r[P1] and r[P3].
	 */
	Op_OffsetLimit = 0x91

	/* Opcode: IfNotZero P1 P2 * * *
	** Synopsis: if r[P1]!=0 then r[P1]--, goto P2
	**
	** Register P1 must contain an integer.  If the content of register P1 is
	** initially greater than zero, then decrement the value in register P1.
	** If it is non-zero (negative or positive) and then also jump to P2.
	** If register P1 is initially zero, leave it unchanged and fall through.
	 */
	Op_IfNotZero = 0x92

	/* Opcode: DecrJumpZero P1 P2 * * *
	 ** Synopsis: if (--r[P1])==0 goto P2
	 **
	 ** Register P1 must hold an integer.  Decrement the value in P1
	 ** and jump to P2 if the new value is exactly zero.
	 */
	Op_DecrJumpZero = 0x93

	/* Opcode: AggStep * P2 P3 P4 P5
	** Synopsis: accum=r[P3] step(r[P2@P5])
	**
	** Execute the xStep function for an aggregate.
	** The function has P5 arguments.  P4 is a pointer to the
	** FuncDef structure that specifies the function.  Register P3 is the
	** accumulator.
	**
	** The P5 arguments are taken from register P2 and its
	** successors.
	 */
	Op_AddgStep = 0x94

	/* Opcode: AggInverse * P2 P3 P4 P5
	 ** Synopsis: accum=r[P3] inverse(r[P2@P5])
	 **
	 ** Execute the xInverse function for an aggregate.
	 ** The function has P5 arguments.  P4 is a pointer to the
	 ** FuncDef structure that specifies the function.  Register P3 is the
	 ** accumulator.
	 **
	 ** The P5 arguments are taken from register P2 and its
	 ** successors.
	 */
	Op_AggInverse = 0x95

	/* Opcode: AggStep1 P1 P2 P3 P4 P5
	** Synopsis: accum=r[P3] step(r[P2@P5])
	**
	** Execute the xStep (if P1==0) or xInverse (if P1!=0) function for an
	** aggregate.  The function has P5 arguments.  P4 is a pointer to the
	** FuncDef structure that specifies the function.  Register P3 is the
	** accumulator.
	**
	** The P5 arguments are taken from register P2 and its
	** successors.
	**
	** This opcode is initially coded as OP_AggStep0.  On first evaluation,
	** the FuncDef stored in P4 is converted into an sqlite3_context and
	** the opcode is changed.  In this way, the initialization of the
	** sqlite3_context only happens once, instead of on each call to the
	** step function.
	 */
	Op_AggStep1 = 0x96

	/* Opcode: AggFinal P1 P2 * P4 *
	 ** Synopsis: accum=r[P1] N=P2
	 **
	 ** P1 is the memory location that is the accumulator for an aggregate
	 ** or window function.  Execute the finalizer function
	 ** for an aggregate and store the result in P1.
	 **
	 ** P2 is the number of arguments that the step function takes and
	 ** P4 is a pointer to the FuncDef for this function.  The P2
	 ** argument is not used by this opcode.  It is only there to disambiguate
	 ** functions that can take varying numbers of arguments.  The
	 ** P4 argument is only needed for the case where
	 ** the step function was not previously called.
	 */
	Op_AggFinal = 0x97

	/* Opcode: AggValue * P2 P3 P4 *
	** Synopsis: r[P3]=value N=P2
	**
	** Invoke the xValue() function and store the result in register P3.
	**
	** P2 is the number of arguments that the step function takes and
	** P4 is a pointer to the FuncDef for this function.  The P2
	** argument is not used by this opcode.  It is only there to disambiguate
	** functions that can take varying numbers of arguments.  The
	** P4 argument is only needed for the case where
	** the step function was not previously called.
	 */
	Op_AggValue = 0x98

	/* Opcode: Checkpoint P1 P2 P3 * *
	**
	** Checkpoint database P1. This is a no-op if P1 is not currently in
	** WAL mode. Parameter P2 is one of SQLITE_CHECKPOINT_PASSIVE, FULL,
	** RESTART, or TRUNCATE.  Write 1 or 0 into mem[P3] if the checkpoint returns
	** SQLITE_BUSY or not, respectively.  Write the number of pages in the
	** WAL after the checkpoint into mem[P3+1] and the number of pages
	** in the WAL that have been checkpointed after the checkpoint
	** completes into mem[P3+2].  However on an error, mem[P3+1] and
	** mem[P3+2] are initialized to -1.
	 */
	Op_Checkpoint = 0x99

	/* Opcode: JournalMode P1 P2 P3 * *
	**
	** Change the journal mode of database P1 to P3. P3 must be one of the
	** PAGER_JOURNALMODE_XXX values. If changing between the various rollback
	** modes (delete, truncate, persist, off and memory), this is a simple
	** operation. No IO is required.
	**
	** If changing into or out of WAL mode the procedure is more complicated.
	**
	** Write a string containing the final journal-mode to register P2.
	 */
	Op_JournalMode = 0x9a

	/* Opcode: Vacuum P1 P2 * * *
	**
	** Vacuum the entire database P1.  P1 is 0 for "main", and 2 or more
	** for an attached database.  The "temp" database may not be vacuumed.
	**
	** If P2 is not zero, then it is a register holding a string which is
	** the file into which the result of vacuum should be written.  When
	** P2 is zero, the vacuum overwrites the original database.
	 */
	Op_Vacuum = 0x9b

	/* Opcode: IncrVacuum P1 P2 * * *
	**
	** Perform a single step of the incremental vacuum procedure on
	** the P1 database. If the vacuum has finished, jump to instruction
	** P2. Otherwise, fall through to the next instruction.
	 */
	Op_IncrVacuum = 0x9c

	/* Opcode: Expire P1 P2 * * *
	**
	** Cause precompiled statements to expire.  When an expired statement
	** is executed using sqlite3_step() it will either automatically
	** reprepare itself (if it was originally created using sqlite3_prepare_v2())
	** or it will fail with SQLITE_SCHEMA.
	**
	** If P1 is 0, then all SQL statements become expired. If P1 is non-zero,
	** then only the currently executing statement is expired.
	**
	** If P2 is 0, then SQL statements are expired immediately.  If P2 is 1,
	** then running SQL statements are allowed to continue to run to completion.
	** The P2==1 case occurs when a CREATE INDEX or similar schema change happens
	** that might help the statement run faster but which does not affect the
	** correctness of operation.
	 */
	Op_Expire = 0x9d

	/* Opcode: CursorLock P1 * * * *
	 **
	 ** Lock the btree to which cursor P1 is pointing so that the btree cannot be
	 ** written by an other cursor.
	 */
	Op_CursorLock = 0x9e

	/* Opcode: CursorUnlock P1 * * * *
	**
	** Unlock the btree to which cursor P1 is pointing so that it can be
	** written by other cursors.
	 */
	Op_CursorUnlock = 0x9f

	/* Opcode: TableLock P1 P2 P3 P4 *
	** Synopsis: iDb=P1 root=P2 write=P3
	**
	** Obtain a lock on a particular table. This instruction is only used when
	** the shared-cache feature is enabled.
	**
	** P1 is the index of the database in sqlite3.aDb[] of the database
	** on which the lock is acquired.  A readlock is obtained if P3==0 or
	** a write lock if P3==1.
	**
	** P2 contains the root-page of the table to lock.
	**
	** P4 contains a pointer to the name of the table being locked. This is only
	** used to generate an error message if the lock cannot be obtained.
	 */
	Op_TableLock = 0xa0

	/* Opcode: VBegin * * * P4 *
	**
	** P4 may be a pointer to an sqlite3_vtab structure. If so, call the
	** xBegin method for that table.
	**
	** Also, whether or not P4 is set, check that this is not being called from
	** within a callback to a virtual table xSync() method. If it is, the error
	** code will be set to SQLITE_LOCKED.
	 */
	Op_VBegin = 0xa1

	/* Opcode: VCreate P1 P2 * * *
	**
	** P2 is a register that holds the name of a virtual table in database
	** P1. Call the xCreate method for that table.
	 */
	Op_VCreate = 0xa2

	/* Opcode: VDestroy P1 * * P4 *
	**
	** P4 is the name of a virtual table in database P1.  Call the xDestroy method
	** of that table.
	 */
	Op_VDestroy = 0xa3

	/* Opcode: VOpen P1 * * P4 *
	**
	** P4 is a pointer to a virtual table object, an sqlite3_vtab structure.
	** P1 is a cursor number.  This opcode opens a cursor to the virtual
	** table and stores that cursor in P1.
	 */
	Op_VOpen = 0xa4

	/* Opcode: VFilter P1 P2 P3 P4 *
	** Synopsis: iplan=r[P3] zplan='P4'
	**
	** P1 is a cursor opened using VOpen.  P2 is an address to jump to if
	** the filtered result set is empty.
	**
	** P4 is either NULL or a string that was generated by the xBestIndex
	** method of the module.  The interpretation of the P4 string is left
	** to the module implementation.
	**
	** This opcode invokes the xFilter method on the virtual table specified
	** by P1.  The integer query plan parameter to xFilter is stored in register
	** P3. Register P3+1 stores the argc parameter to be passed to the
	** xFilter method. Registers P3+2..P3+1+argc are the argc
	** additional parameters which are passed to
	** xFilter as argv. Register P3+2 becomes argv[0] when passed to xFilter.
	**
	** A jump is made to P2 if the result set after filtering would be empty.
	 */
	Op_VFilter = 0xa5

	/* Opcode: VColumn P1 P2 P3 * P5
	** Synopsis: r[P3]=vcolumn(P2)
	**
	** Store in register P3 the value of the P2-th column of
	** the current row of the virtual-table of cursor P1.
	**
	** If the VColumn opcode is being used to fetch the value of
	** an unchanging column during an UPDATE operation, then the P5
	** value is OPFLAG_NOCHNG.  This will cause the sqlite3_vtab_nochange()
	** function to return true inside the xColumn method of the virtual
	** table implementation.  The P5 column might also contain other
	** bits (OPFLAG_LENGTHARG or OPFLAG_TYPEOFARG) but those bits are
	** unused by OP_VColumn.
	 */
	Op_VColumn = 0xa6

	/* Opcode: VNext P1 P2 * * *
	**
	** Advance virtual table P1 to the next row in its result set and
	** jump to instruction P2.  Or, if the virtual table has reached
	** the end of its result set, then fall through to the next instruction.
	 */
	Op_VNext = 0xa7

	/* Opcode: VRename P1 * * P4 *
	**
	** P4 is a pointer to a virtual table object, an sqlite3_vtab structure.
	** This opcode invokes the corresponding xRename method. The value
	** in register P1 is passed as the zName argument to the xRename method.
	 */
	Op_VRename = 0xa8

	/* Opcode: VUpdate P1 P2 P3 P4 P5
	** Synopsis: data=r[P3@P2]
	**
	** P4 is a pointer to a virtual table object, an sqlite3_vtab structure.
	** This opcode invokes the corresponding xUpdate method. P2 values
	** are contiguous memory cells starting at P3 to pass to the xUpdate
	** invocation. The value in register (P3+P2-1) corresponds to the
	** p2th element of the argv array passed to xUpdate.
	**
	** The xUpdate method will do a DELETE or an INSERT or both.
	** The argv[0] element (which corresponds to memory cell P3)
	** is the rowid of a row to delete.  If argv[0] is NULL then no
	** deletion occurs.  The argv[1] element is the rowid of the new
	** row.  This can be NULL to have the virtual table select the new
	** rowid for itself.  The subsequent elements in the array are
	** the values of columns in the new row.
	**
	** If P2==1 then no insert is performed.  argv[0] is the rowid of
	** a row to delete.
	**
	** P1 is a boolean flag. If it is set to true and the xUpdate call
	** is successful, then the value returned by sqlite3_last_insert_rowid()
	** is set to the value of the rowid for the row just inserted.
	**
	** P5 is the error actions (OE_Replace, OE_Fail, OE_Ignore, etc) to
	** apply in the case of a constraint failure on an insert or update.
	 */
	Op_VUpdate = 0xa9

	/* Opcode: Pagecount P1 P2 * * *
	**
	** Write the current number of pages in database P1 to memory cell P2.
	 */
	Op_Pagecount = 0xaa

	/* Opcode: MaxPgcnt P1 P2 P3 * *
	**
	** Try to set the maximum page count for database P1 to the value in P3.
	** Do not let the maximum page count fall below the current page count and
	** do not change the maximum page count value if P3==0.
	**
	** Store the maximum page count after the change in register P2.
	 */
	Op_MaxPgcnt = 0xab

	/* Opcode: Function P1 P2 P3 P4 *
	** Synopsis: r[P3]=func(r[P2@NP])
	**
	** Invoke a user function (P4 is a pointer to an sqlite3_context object that
	** contains a pointer to the function to be run) with arguments taken
	** from register P2 and successors.  The number of arguments is in
	** the sqlite3_context object that P4 points to.
	** The result of the function is stored
	** in register P3.  Register P3 must not be one of the function inputs.
	**
	** P1 is a 32-bit bitmask indicating whether or not each argument to the
	** function was determined to be constant at compile time. If the first
	** argument was constant then bit 0 of P1 is set. This is used to determine
	** whether meta data associated with a user function argument using the
	** sqlite3_set_auxdata() API may be safely retained until the next
	** invocation of this opcode.
	**
	** See also: AggStep, AggFinal, PureFunc
	 */
	Op_Function = 0xac

	/* Opcode: PureFunc P1 P2 P3 P4 *
	** Synopsis: r[P3]=func(r[P2@NP])
	**
	** Invoke a user function (P4 is a pointer to an sqlite3_context object that
	** contains a pointer to the function to be run) with arguments taken
	** from register P2 and successors.  The number of arguments is in
	** the sqlite3_context object that P4 points to.
	** The result of the function is stored
	** in register P3.  Register P3 must not be one of the function inputs.
	**
	** P1 is a 32-bit bitmask indicating whether or not each argument to the
	** function was determined to be constant at compile time. If the first
	** argument was constant then bit 0 of P1 is set. This is used to determine
	** whether meta data associated with a user function argument using the
	** sqlite3_set_auxdata() API may be safely retained until the next
	** invocation of this opcode.
	**
	** This opcode works exactly like OP_Function.  The only difference is in
	** its name.  This opcode is used in places where the function must be
	** purely non-deterministic.  Some built-in date/time functions can be
	** either determinitic of non-deterministic, depending on their arguments.
	** When those function are used in a non-deterministic way, they will check
	** to see if they were called using OP_PureFunc instead of OP_Function, and
	** if they were, they throw an error.
	**
	** See also: AggStep, AggFinal, Function
	 */
	Op_PureFunc = 0xad

	/* Opcode: Trace P1 P2 * P4 *
	**
	** Write P4 on the statement trace output if statement tracing is
	** enabled.
	**
	** Operand P1 must be 0x7fffffff and P2 must positive.
	 */
	Op_Trace = 0xae

	/* Opcode: Init P1 P2 P3 P4 *
	** Synopsis: Start at P2
	**
	** Programs contain a single instance of this opcode as the very first
	** opcode.
	**
	** If tracing is enabled (by the sqlite3_trace()) interface, then
	** the UTF-8 string contained in P4 is emitted on the trace callback.
	** Or if P4 is blank, use the string returned by sqlite3_sql().
	**
	** If P2 is not zero, jump to instruction P2.
	**
	** Increment the value of P1 so that OP_Once opcodes will jump the
	** first time they are evaluated for this run.
	**
	** If P3 is not zero, then it is an address to jump to if an SQLITE_CORRUPT
	** error is encountered.
	 */
	Op_Init = 0xaf

	/* Opcode: CursorHint P1 * * P4 *
	**
	** Provide a hint to cursor P1 that it only needs to return rows that
	** satisfy the Expr in P4.  TK_REGISTER terms in the P4 expression refer
	** to values currently held in registers.  TK_COLUMN terms in the P4
	** expression refer to columns in the b-tree to which cursor P1 is pointing.
	 */
	Op_CursorHint = 0xb0

	/* Opcode:  Abortable   * * * * *
	**
	** Verify that an Abort can happen.  Assert if an Abort at this point
	** might cause database corruption.  This opcode only appears in debugging
	** builds.
	**
	** An Abort is safe if either there have been no writes, or if there is
	** an active statement journal.
	 */
	Op_Abortable = 0xb1

	/* Opcode:  ReleaseReg   P1 P2 P3 * P5
	** Synopsis: release r[P1@P2] mask P3
	**
	** Release registers from service.  Any content that was in the
	** the registers is unreliable after this opcode completes.
	**
	** The registers released will be the P2 registers starting at P1,
	** except if bit ii of P3 set, then do not release register P1+ii.
	** In other words, P3 is a mask of registers to preserve.
	**
	** Releasing a register clears the Mem.pScopyFrom pointer.  That means
	** that if the content of the released register was set using OP_SCopy,
	** a change to the value of the source register for the OP_SCopy will no longer
	** generate an assertion fault in sqlite3VdbeMemAboutToChange().
	**
	** If P5 is set, then all released registers have their type set
	** to MEM_Undefined so that any subsequent attempt to read the released
	** register (before it is reinitialized) will generate an assertion fault.
	**
	** P5 ought to be set on every call to this opcode.
	** However, there are places in the code generator will release registers
	** before their are used, under the (valid) assumption that the registers
	** will not be reallocated for some other purpose before they are used and
	** hence are safe to release.
	**
	** This opcode is only available in testing and debugging builds.  It is
	** not generated for release builds.  The purpose of this opcode is to help
	** validate the generated bytecode.  This opcode does not actually contribute
	** to computing an answer.
	 */
	OP_ReleaseReg = 0xb2
)
