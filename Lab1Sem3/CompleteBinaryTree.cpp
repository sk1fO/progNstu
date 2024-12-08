#include <iostream>
#include <vector>
#include <queue>

using namespace std;

class TreeNode {
public:
    int value;
    TreeNode *left, *right;
    TreeNode(int val) : value(val), left(nullptr), right(nullptr) {}
};

class CompleteBinaryTree {
private:
    TreeNode *root;
    vector<TreeNode*> nodes;

public:
    CompleteBinaryTree() : root(nullptr) {}

    void insert(int value) {
        TreeNode *new_node = new TreeNode(value);
        if (!root) {
            root = new_node;
            nodes.push_back(new_node);
            return;
        }

        for (auto node : nodes) {
            if (!node-&gt;left) {
                node-&gt;left = new_node;
                break;
            } else if (!node-&gt;right) {
                node-&gt;right = new_node;
                break;
            }
        }

        nodes.push_back(new_node);
    }

    void level_order() {
        if (!root) return;
        queue<treenode*> q;
        q.push(root);
        while (!q.empty()) {
            TreeNode *node = q.front();
            q.pop();
            cout &lt;&lt; node-&gt;value &lt;&lt; " ";
            if (node-&gt;left) q.push(node-&gt;left);
            if (node-&gt;right) q.push(node-&gt;right);
        }
    }
};

int main() {
    CompleteBinaryTree tree;
    tree.insert(1);
    tree.insert(2);
    tree.insert(3);
    tree.insert(4);
    tree.insert(5);
    tree.insert(6);
    tree.level_order();  // Output: 1 2 3 4 5 6

    return 0;
}