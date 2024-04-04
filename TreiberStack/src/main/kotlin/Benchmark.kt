package concurrentStack

import concurrentStack.stack.ConcurrentEliminationStack
import concurrentStack.stack.ConcurrentStack
import concurrentStack.stack.StackInterface
import kotlinx.coroutines.launch
import kotlinx.coroutines.runBlocking
import kotlin.system.measureTimeMillis

class Benchmark {
    private fun pushMany(stack: StackInterface<Int>, n: Int = 1000) {
        for (i in 1..n) {
            stack.push(i)
        }
    }

    private fun popMany(stack: StackInterface<Int>, n: Int = 1000) {
        for (i in 1..n) {
            try {
                stack.pop()
            } catch (e: Exception) {
            }
        }
    }

    fun benchAlternate(
        nruns: Int = 5,
        noperations: Int = 1000,
        njobs: Int = 2,
        waitSteps: Int = 100
    ): Pair<Long, Long> {
        val concurrentStack = ConcurrentStack<Int>()
        val eliminationStack = ConcurrentEliminationStack<Int>(waitSteps = waitSteps)
        val t1s = Array<Long>(nruns) { 0 }
        val t2s = Array<Long>(nruns) { 0 }

        for (i in 0..<nruns) {
            val t: Long = runBlocking {
                measureTimeMillis {
                    for (j in 0..<njobs) {
                        if (j % 2 == 0) {
                            launch { pushMany(concurrentStack, noperations) }.join()
                        } else {
                            launch { popMany(concurrentStack, noperations) }.join()
                        }
                    }
                }
            }
            t1s[i] = t
        }

        for (i in 0..<nruns) {
            val t: Long = runBlocking {
                measureTimeMillis {
                    for (j in 0..<njobs) {
                        if (j % 2 == 0) {
                            launch { pushMany(eliminationStack, noperations) }.join()
                        } else {
                            launch { popMany(eliminationStack, noperations) }.join()
                        }
                    }
                }
            }
            t2s[i] = t
        }

        return Pair(t1s.sum() / nruns, t2s.sum() / nruns)
    }

    fun benchPushOnly(
        nruns: Int = 5,
        noperations: Int = 1000,
        njobs: Int = 2,
        waitSteps: Int = 100
    ): Pair<Long, Long> {
        val concurrentStack = ConcurrentStack<Int>()
        val eliminationStack = ConcurrentEliminationStack<Int>(waitSteps = waitSteps)
        val t1s = Array<Long>(nruns) { 0 }
        val t2s = Array<Long>(nruns) { 0 }

        for (i in 0..<nruns) {
            val t: Long = runBlocking {
                measureTimeMillis {
                    for (j in 0..<njobs) {
                        if (j % 2 == 0) {
                            launch { pushMany(concurrentStack, noperations) }.join()
                        } else {
                            launch { pushMany(concurrentStack, noperations) }.join()
                        }
                    }
                }
            }
            t1s[i] = t
        }

        for (i in 0..<nruns) {
            val t: Long = runBlocking {
                measureTimeMillis {
                    for (j in 0..<njobs) {
                        if (j % 2 == 0) {
                            launch { pushMany(eliminationStack, noperations) }.join()
                        } else {
                            launch { pushMany(eliminationStack, noperations) }.join()
                        }
                    }
                }
            }
            t2s[i] = t
        }

        return Pair(t1s.sum() / nruns, t2s.sum() / nruns)
    }

    fun benchPopOnly(
        nruns: Int = 5,
        noperations: Int = 1000,
        njobs: Int = 2,
        waitSteps: Int = 100
    ): Pair<Long, Long> {
        val concurrentStack = ConcurrentStack<Int>()
        val eliminationStack = ConcurrentEliminationStack<Int>(waitSteps = waitSteps)
        val t1s = Array<Long>(nruns) { 0 }
        val t2s = Array<Long>(nruns) { 0 }

        for (i in 0..<nruns) {
            val t: Long = runBlocking {
                measureTimeMillis {
                    for (j in 0..<njobs) {
                        if (j % 2 == 0) {
                            launch { popMany(concurrentStack, noperations) }.join()
                        } else {
                            launch { popMany(concurrentStack, noperations) }.join()
                        }
                    }
                }
            }
            t1s[i] = t
        }

        for (i in 0..<nruns) {
            val t: Long = runBlocking {
                measureTimeMillis {
                    for (j in 0..<njobs) {
                        if (j % 2 == 0) {
                            launch { popMany(eliminationStack, noperations) }.join()
                        } else {
                            launch { popMany(eliminationStack, noperations) }.join()
                        }
                    }
                }
            }
            t2s[i] = t
        }

        return Pair(t1s.sum() / nruns, t2s.sum() / nruns)
    }

}