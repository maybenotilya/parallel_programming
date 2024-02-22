package concurrentStack.stack

class StackNode<T>(val value: T, val next: StackNode<T>? = null) {}