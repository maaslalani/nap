# Python3 implementation of QuickSort
# Function to find the partition position
def partition(array, low, high):
  
    # Choose the rightmost element as pivot
    pivot = array[high]
  
    # Pointer for greater element
    i = low - 1
  
    # Traverse through all elements
    # compare each element with pivot
    for j in range(low, high):
        if array[j] <= pivot:
            # If element smaller than pivot is found
            # swap it with the greater element pointed by i
            i = i + 1
  
            # Swapping element at i with element at j
            (array[i], array[j]) = (array[j], array[i])
  
    # Swap the pivot element with 
    # e greater element specified by i
    (array[i + 1], array[high]) = (array[high], array[i + 1])
  
    # Return the position from where partition is done
    return i + 1
  
# Function to perform quicksort
  
  
def quick_sort(array, low, high):
    if low < high:
  
        # Find pivot element such that
        # element smaller than pivot are on the left
        # element greater than pivot are on the right
        pi = partition(array, low, high)
  
        # Recursive call on the left of pivot
        quick_sort(array, low, pi - 1)
  
        # Recursive call on the right of pivot
        quick_sort(array, pi + 1, high)
  
  
# Driver code
array = [10, 7, 8, 9, 1, 5]
quick_sort(array, 0, len(array) - 1)
  
print(f'Sorted array: {array}')
