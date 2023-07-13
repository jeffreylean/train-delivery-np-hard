## Bigpay assessment - Trains

The problem in the assessment is a np hard problem which falls into the same category as TSP (Travel Sales Person). I end up wrote two solution for this, both of them are
heuristic approach to the problem, and might not be the most optimise solution.

### Solution 1

Solution 1 is a much more simpler approach to the problem, as I used a greedy algorithm approach to solve the problem by using djikstra algorithm to calculate the
shortest distance between the train and the next destination.

For single train only

1. First identify which are the closest package to the train, here djikstra is used, after that the train will go and pickup the package.
2. If there are multiple package available, the next step will be either deliver the pick up package to the destination or go and pick up another package.
   The deciding factor is which one is the closest, the package destination? or the next package location.
3. 1 & 2 is repeat until all package has been deliver

For multiple train

1. Assign the "best" package to each train respectively. The "best" here take into the consideration of weight and distance of a package. If there's a package which
   is too far, then it might be the best that the train that have the capacity prioritize that first, where train with lesser capacity can focus on the lighter one.
   the priority is calculated by weight/distance ratio.
2. After the assignment has been done, each train will start to pick up the package.
3. Each train then will have to go thru the continue picking up or directly delivery current package to the destination decision. This is decide by the shorter distance.4. 1,2 & 3 is repeat until there's no package left to deliver.

#### To run the program

- Go to the solution1 directory and run `go run main.go`
- The input are in `example.txt` file, you can modify the input and rerun the program to see the changes.

### Solution 2

Solution 2 is a much more optimal approach as I used an algorithm specificly to solve this type of problem which is simulated annealing. Stimulate annealing is a
probabilistic optimization inspired by the process of annealing in metallurgy, where a material is heated and slowly cooled to reduce the defect and increase the
size of its crystal structure. In computer science, this process is mimics as an approach to find a good approximation to the global minimum of a given
function in a large search space. Think of the hill climbing, which move by finding better neighbour, but the outcome might easily be just a local optimum, simulated annealing use random point instead.

The core of simulated annealing is the equation below:

Prob(To reach global optimum) ~ 1 - exp(delta E / T))

This is an equation borrowed from thermal physics, where the acception probability the delta of E (energy), which is E2 - E1 divided by T (temperature). In our case,
the energy would be our time taken for all package to be delivered. The higher the T, the larger the probability to be accept, so it is used as control parameter that
influence the algorithm's exploration. Means that if in this iteration, if the time taken is longer, but if the T is high enough, the probability of we accepting it
still high. By accepting the worse solution, it helps us excaping the local optima and exporing more search space in the finding of global optima.

The idea is that the initiation of the simulation, the T value should be high enough in the exploratory phase. As the algorithm progress, T will decrease. As T get
smaller, the algorithm become less likely to accept worse solution and more focused on exploiting the promising area where the search space it has discovered.

There's also the concept of neighbour which refer to a state that is similar to the current state but with some small change applied. The neighbour energy (time taken) will be used to compare with the current state, to see if it is a better solution or worst.

So how it works with this problem, here's the iteration.

- create initial state for all the train. Meaning assignment than any package available with respect to the capacity. Here I use the same concept in solution 1 where I
  assign the closest package to the train. Also the initial state involve path as well, so here we will be using the shortest path to the package assigned.
- for n interation
  - Create a neighbour by making small changes to the current state, there are 2 way of doing that:
    1. Get 2 random train (if there are multiple train), assign package's of train1 to train2
    2. Get one random train, and swap the order of the package, so that the route of the train will be recalculated and change.
  - If the time taken is shorter (better route):
    - accept this
  - Else:
    - Calculate the acceptance probability.
    - Generate a random number
    - if random number > acceptance probability:
      - Reject
    - Else:
      - Accept

After all iteration, we should get a optima energy, which is our optimal time taken to deliver all the package.

#### To run the program

- Go to the solution1 directory and run `go run main.go`
- The input are in `example.txt` file, you can modify the input and rerun the program to see the changes.
