package concurrentStack

import concurrentStack.stack.ConcurrentStack
import concurrentStack.stack.IntStack
import org.jetbrains.kotlinx.lincheck.annotations.Operation
import org.jetbrains.kotlinx.lincheck.check
import org.jetbrains.kotlinx.lincheck.strategy.managed.modelchecking.ModelCheckingOptions
import org.jetbrains.kotlinx.lincheck.strategy.stress.StressOptions
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.assertThrows
import java.util.EmptyStackException

class ConcurrentStackTest {

    @Test
    fun `Push-pop`() {
        val stack = ConcurrentStack<Int>()
        for (i in 1..10) {
            stack.push(i)
        }

        for (i in 1..5) {
            stack.pop()
        }
        assertEquals(stack.pop(), 5)
        assertEquals(stack.top(), 4)

    }

    @Test
    fun `Push-pop equal elements`() {
        val stack = ConcurrentStack<Int>()
        for (i in 1..10) {
            stack.push(1)
        }
        for (i in 1..9) {
            stack.pop()
        }

        assertEquals(stack.top(), 1)
        assertEquals(stack.pop(), 1)
        assertThrows<EmptyStackException> { stack.pop() }
    }

    @Test
    fun `Pop empty`() {
        val stack = ConcurrentStack<Int>()
        assertEquals(stack.top(), null)
        assertThrows<EmptyStackException> { stack.pop() }
    }
}

class ConcurrentStackStressTest {
    private val stack = ConcurrentStack<Int>()

    @Operation
    fun push(x: Int) = stack.push(x)

    @Operation
    fun pop() = stack.pop()

    @Operation
    fun top() = stack.top()

    @Test
    fun stressTest() = StressOptions()
        .sequentialSpecification(IntStack::class.java)
        .check(this::class.java)

    @Test
    fun modelTest() = ModelCheckingOptions()
        .checkObstructionFreedom()
        .sequentialSpecification(IntStack::class.java)
        .check(this::class.java)
}
