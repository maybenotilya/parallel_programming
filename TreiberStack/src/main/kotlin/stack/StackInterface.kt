package concurrentStack.stack

interface StackInterface<T> {
    fun push(x: T) : Unit
    fun pop() : T
    fun top() : T?
}