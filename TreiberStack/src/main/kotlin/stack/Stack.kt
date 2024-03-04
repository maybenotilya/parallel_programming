package concurrentStack.stack

import java.util.EmptyStackException

class IntStack<T> : StackInterface<Int> {
    private val stack = mutableListOf<Int>()

    override fun push(x: Int) {
        stack.add(x)
    }

    override fun pop(): Int {
        return stack.removeLastOrNull() ?: throw EmptyStackException()
    }

    override fun top(): Int? = stack.lastOrNull()
}