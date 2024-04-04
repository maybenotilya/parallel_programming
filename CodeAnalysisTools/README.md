# Задача 2. Инструменты анализа кода

Возьмите любой проект с открытым исходным кодом на C/C++ без элементов OpenMP,
содержащий элементы параллельного программирования. Проанализируйте его
(насколько это возможно, всесторонне) с помощью инструментов Helgrind и
ThreadSanitizer или их аналогом.
Какие предупреждения выдаёт инструмент? Указывают ли они на реальную проблему
или хотя бы на ту, которая потенциально может возникнуть?
Далее необходимо внести в проект некоторую гонку данных (не совсем уж
искусственную) и обнаружить её с помощью инструмента.

Комментарий. Язык программирования — C/С++ в силу особенностей инструмента

# Game of Life

В [репозитории](https://github.com/maybenotilya/game-of-life) реализована многопоточная версия "Игры в жизнь". Запуск программы поочередно симулирует игру, используя от 1 до введенного максимального количества потоков. Многопоточная работа ведется построчно: каждому потоку выделяется равное (или примерно равное в случае, когда высота не делится на количество потоков) количество строк

Запускается следующим образом

```
./gol_parallel <input file> <num steps> <output file> <max number of threads>
```

# Helgrind

Анализ запускался со следующими параметрами:

```
valgrind --tool=helgrind ./gol_parallel random.txt 100 output.txt 8
```

Результат:
```
ERROR SUMMARY: 280007 errors from 4 contexts (suppressed: 0 from 0)
```
Рассмотрим, на что он ругается:

### Thread #1: pthread_barrier_init: barrier is already initialised

```
Thread #1: pthread_barrier_init: barrier is already initialised
at 0x4853494: pthread_barrier_init (hg_intercepts.c:1873)
by 0x1099A5: walltime_of_threads (in /home/maybenotilya/proj/game-of-life/src/gol_parallel)
by 0x109DFD: main (in /home/maybenotilya/proj/game-of-life/src/gol_parallel)
```

Причиной возникновения является то, что барьер объявлен в глобальной области видимости, а инициализируется в функции `walltime_of_threads`, которая многократно вызывается при работе программы. Это приводит к многократной инициализации одного и того же барьера. Исправить это можно создавая в функции `walltime_of_threads` каждый раз новый барьер и передавая его в качестве аргумента для `entry_function`, которая отвечает за работу с потоками. Однако кажется это не сильно критическая ошибка.

### Possible data race during write of size 1 at 0x4B06961 by thread #4

```
Possible data race during write of size 1 at 0x4B06961 by thread #4
Locks held: none
   at 0x109801: update_matrix (in /home/maybenotilya/proj/game-of-life/src/gol_parallel)
   by 0x1098A7: entry_function (in /home/maybenotilya/proj/game-of-life/src/gol_parallel)
   by 0x4851B3B: mythread_wrapper (hg_intercepts.c:406)
   by 0x490E559: start_thread (pthread_create.c:447)
   by 0x498B873: clone (clone.S:100)

This conflicts with a previous write of size 1 by thread #3
Locks held: none
   at 0x109801: update_matrix (in /home/maybenotilya/proj/game-of-life/src/gol_parallel)
   by 0x1098A7: entry_function (in /home/maybenotilya/proj/game-of-life/src/gol_parallel)
   by 0x4851B3B: mythread_wrapper (hg_intercepts.c:406)
   by 0x490E559: start_thread (pthread_create.c:447)
   by 0x498B873: clone (clone.S:100)
 Address 0x4b06961 is 1 bytes inside a block of size 102 alloc'd
   at 0x484C993: calloc (vg_replace_malloc.c:1595)
   by 0x1092F4: allocate_matrix (in /home/maybenotilya/proj/game-of-life/src/gol_parallel)
   by 0x109D8A: main (in /home/maybenotilya/proj/game-of-life/src/gol_parallel)
 Block was alloc'd by thread #1
```

Гонка данных возникает в функции `update_matrix`. Причиной этому является то, что построчная обработка идет от `begin_row` до `end_row` включительно. Выбор строк матрицы происходит в функции `entry_function` на 116-118 строках:

```
int bound = height / n_threads;
int begin_row = i_thread * bound;
int end_row = begin_row + bound;
```

Предположим, что `bound` получился равным `n`. Тогда потока с идентификатором `i_thread = 0` получит в пользование строки с `0` по `n` включительно. Однако, поток с `i_thread = 1` получит строки с `n` до `2n` включительно. Видно, что оба потока разделяют строку с номером `n`. 

После изменений:

```
Строка 94
-    for(i=begin_row;i<end_row+1;i++){
+    for(i=begin_row;i<end_row;i++){
```

```
Строка 122
-    if(i_thread==n_threads-1) end_row=height;
+    if(i_thread==n_threads-1) end_row=height + 1;
```

данная ошибка пропала. Принцип работы алгоритма не изменился.

# ThreadSanitizer

При запуске с флагами компиляции `-fsanitize=thread -g` возникла следующая ошибка:

```
FATAL: ThreadSanitizer: unexpected memory mapping 0x64e316430000-0x64e316431000
```

Помогло выключение ASLR в `/proc/sys/kernel/randomize_va_space` (переключением флага с `2` на `0`).

На версии до исправлений ошибок Helgrind'а выдавались предупреждения следующего вида:

```
WARNING: ThreadSanitizer: data race (pid=23153)
  Write of size 1 at 0x7b1c0004c4b8 by thread T9:
    #0 update_matrix /home/maybenotilya/proj/game-of-life/src/gol_parallel.c:99 (gol_parallel+0x1a96) (BuildId: f1479316e208041f350e071aad0ed5aa469333be)
    #1 entry_function /home/maybenotilya/proj/game-of-life/src/gol_parallel.c:127 (gol_parallel+0x1c92) (BuildId: f1479316e208041f350e071aad0ed5aa469333be)

  Previous write of size 1 at 0x7b1c0004c4b8 by thread T10:
    #0 update_matrix /home/maybenotilya/proj/game-of-life/src/gol_parallel.c:99 (gol_parallel+0x1a96) (BuildId: f1479316e208041f350e071aad0ed5aa469333be)
    #1 entry_function /home/maybenotilya/proj/game-of-life/src/gol_parallel.c:127 (gol_parallel+0x1c92) (BuildId: f1479316e208041f350e071aad0ed5aa469333be)
```

Очевидно, это та же проблема с гонкой данных, описанная выше. После исправлений все предупреждения пропали.