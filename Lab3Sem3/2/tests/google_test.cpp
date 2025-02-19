#include <gtest/gtest.h>
#include "../src/hashtable.h"

TEST(HashTableTest, InsertAndSearch) {
    HashTable ht;
    ht.insert("apple", "banana");
    EXPECT_EQ(ht.search("apple"), "banana");
    ht.insert("apple", "cherry");
    EXPECT_EQ(ht.search("apple"), "cherry");
}

TEST(HashTableTest, Remove) {
    HashTable ht;
    ht.insert("apple", "banana");
    ht.insert("orange", "grape");
    ht.remove("apple");
    EXPECT_EQ(ht.search("apple"), "not exists");
    EXPECT_EQ(ht.search("orange"), "grape");
}

TEST(HashTableTest, Print) {
    HashTable ht;
    ht.insert("apple", "banana");
    ht.insert("orange", "grape");
    std::stringstream ss;
    std::streambuf* old_cout = std::cout.rdbuf(ss.rdbuf());
    ht.printTable();
    std::cout.rdbuf(old_cout);
    EXPECT_NE(ss.str().find("[apple:banana]"), std::string::npos);
    EXPECT_NE(ss.str().find("[orange:grape]"), std::string::npos);
}