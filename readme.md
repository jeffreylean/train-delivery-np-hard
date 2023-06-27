## Bigpay assessment - Trains

The problem in the assessment is a np hard problem which falls into the same category as TSP (Travel Sales Person). I end up wrote two solution for this, both of them are
heuristic approach to the problem, and might not be the most optimise solution.

### Solution 1

Solution 1 is a much more simpler approach to the problem, as I used a greedy algorithm approach to solve the problem by using djikstra algorithm to calculate the
shortest distance between the train and the destination. The shortest distance calculated might be a local minima.

1. First identify which are the closest package to the train, here djikstra is used, after that the train will go and pickup the package.
2. If there are multiple package available, the next step will be either deliver the pick up package to the destination or go and pick up another package.
   The deciding factor is which one is the closest, the package destination? or the next package location.
3. 1 & 2 is repeat until all package has been deliver

### Solution 2

Solution 2 is a much more optimal approach as I used an algorithm specificly to solve this type of problem which is Stimulate annealing. Stimulate annealing is a
probabilistic optimization inspired by the process of annealing in metallurgy, where a material is heated and slowly cooled to reduce the defect and increase the
size of its crystal structure. In computer science, this process is mimick as an approach to find a good approximation to the global minimum of a given
function in a large search space.
