package concurrentStack.stack

import kotlinx.atomicfu.atomic
import kotlinx.atomicfu.atomicArrayOfNulls
import java.util.EmptyStackException
import kotlin.random.Random


class Wrapper<T>(val value: T)

class ConcurrentEliminationStack<T>(
    private val eliminationSize: Int = 10, private val waitSteps: Int = 1000
) : StackInterface<T> {
    private val head = atomic<StackNode<T>?>(null)
    private val eliminationArray = atomicArrayOfNulls<Wrapper<T>?>(eliminationSize)

    private fun tryPush(x: T): Boolean {
        val tempHead = head.value
        val newHead = StackNode(x, tempHead)
        return head.compareAndSet(tempHead, newHead)
    }

    private fun tryPop(): T? {
        val tempHead = head.value ?: throw EmptyStackException()
        if (head.compareAndSet(tempHead, tempHead.next)) {
            return tempHead.value
        }
        return null
    }

    private fun pushNoElimination(x: T) {
        while (true) {
            val tempHead = head.value
            val newHead = StackNode(x, tempHead)
            if (head.compareAndSet(tempHead, newHead)) {
                return
            }
        }
    }

    private fun popNoElimination(): T {
        while (true) {
            val tempHead = head.value ?: throw EmptyStackException()
            if (head.compareAndSet(tempHead, tempHead.next)) {
                return tempHead.value
            }
        }
    }

    override fun push(x: T) {
        if (tryPush(x)) {
            return
        }
        val newWrapper = Wrapper(x)
        val pos = Random.nextInt(eliminationSize)
        val randElement = eliminationArray[pos].value
        if (randElement == null) {
            if (eliminationArray[pos].compareAndSet(null, newWrapper)) {
                repeat(waitSteps) {}
                if (eliminationArray[pos].compareAndSet(newWrapper, null)) {
                    pushNoElimination(x)
                }
            } else {
                pushNoElimination(x)
            }
        } else {
            pushNoElimination(x)
        }
    }

    override fun pop(): T {
        val value = tryPop()
        if (value != null) {
            return value
        }
        val pos = Random.nextInt(eliminationSize)
        var step = 0
        while (step < waitSteps) {
            val randElement = eliminationArray[pos].value
            if (randElement != null) {
                if (eliminationArray[pos].compareAndSet(randElement, null)) {
                    return randElement.value
                }
            }
            ++step
        }
        return popNoElimination()
    }

    override fun top(): T? {
        val tempHead = head.value ?: return null
        return tempHead.value
    }
}