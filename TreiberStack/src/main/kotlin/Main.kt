package concurrentStack

fun main() {
    val nruns = 5
    val noperations = 5_000_000
    val waitSteps = 1000
    val njobs = 2
    val bench = Benchmark()

    println("Running push only case:")
    val pushOnlyResult = bench.benchPushOnly(nruns, noperations, njobs, waitSteps)
    println(" - Concurrent stack:  ${pushOnlyResult.first} ms.")
    println(" - Elimination stack: ${pushOnlyResult.second} ms.")
    println()

    println("Running pop only case:")
    val popOnlyResult = bench.benchPopOnly(nruns, noperations, njobs, waitSteps)
    println(" - Concurrent stack:  ${popOnlyResult.first} ms.")
    println(" - Elimination stack: ${popOnlyResult.second} ms.")
    println()

    println("Running alternate case:")
    val alternateResult = bench.benchAlternate(nruns, noperations, njobs, waitSteps)
    println(" - Concurrent stack:  ${alternateResult.first} ms.")
    println(" - Elimination stack: ${alternateResult.second} ms.")
    println()
}