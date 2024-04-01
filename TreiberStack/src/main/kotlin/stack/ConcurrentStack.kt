package concurrentStack.stack

import kotlinx.atomicfu.atomic
import java.util.EmptyStackException

class ConcurrentStack<T> : StackInterface<T> {
    private val head = atomic<StackNode<T>?>(null)

    override fun push(x: T) {
        while (true) {
            val tempHead = head.value
            val newHead = StackNode(x, tempHead)
            if (head.compareAndSet(tempHead, newHead)) {
                return
            }
        }
    }

    override fun pop(): T {
        while (true) {
            val tempHead = head.value ?: throw EmptyStackException()
            if (head.compareAndSet(tempHead, tempHead.next)) {
                return tempHead.value
            }
        }
    }

    override fun top(): T? {
        val tempHead = head.value ?: return null
        return tempHead.value
    }
}