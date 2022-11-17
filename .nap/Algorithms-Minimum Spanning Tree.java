// Java implementation of MST
class GfG {
     // returns the minimum cost of the spanning tree for the required graph
    static int getMinCost(int Vertices) {
        int cost = 0;
        // Calculating cost of MST
        cost = (Vertices * Vertices) - Vertices + 1;
        return cost;
    }
}
 
public static void main(String[] args) {
    int V = 5;
    System.out.println(getMinCost(V));
}