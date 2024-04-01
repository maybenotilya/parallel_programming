buildscript {
    repositories {
        mavenCentral()
    }

    dependencies {
      classpath("org.jetbrains.kotlinx:atomicfu-gradle-plugin:0.23.2")
      classpath("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.8.1-Beta")
    }
}

apply(plugin = "kotlinx-atomicfu")

plugins {
    kotlin("jvm") version "1.9.22"
    id("application")
}

group = "concurrentStack"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    testImplementation("org.junit.jupiter:junit-jupiter:5.9.2")
    testImplementation("org.jetbrains.kotlinx:lincheck:2.25")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.8.1-Beta")
}

tasks.test {
    useJUnitPlatform()
    maxHeapSize = "4g"
}
kotlin {
    jvmToolchain(21)
}

application {
    mainClass = "concurrentStack.MainKt"
}