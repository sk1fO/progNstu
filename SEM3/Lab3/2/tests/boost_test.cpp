#define BOOST_TEST_MODULE HashTableTest
#include <boost/test/unit_test.hpp>
#include "../src/hashtable.h"

BOOST_AUTO_TEST_CASE(test_insert_and_search) {
    HashTable ht;
    ht.insert("apple", "banana");
    BOOST_CHECK_EQUAL(ht.search("apple"), "banana");
    ht.insert("apple", "cherry");
    BOOST_CHECK_EQUAL(ht.search("apple"), "cherry");
}

BOOST_AUTO_TEST_CASE(test_remove) {
    HashTable ht;
    ht.insert("apple", "banana");
    ht.insert("orange", "grape");
    ht.remove("apple");
    BOOST_CHECK_EQUAL(ht.search("apple"), "not exists");
    BOOST_CHECK_EQUAL(ht.search("orange"), "grape");
}

BOOST_AUTO_TEST_CASE(test_print) {
    HashTable ht;
    ht.insert("apple", "banana");
    ht.insert("orange", "grape");
    std::stringstream ss;
    std::streambuf* old_cout = std::cout.rdbuf(ss.rdbuf());
    ht.printTable();
    std::cout.rdbuf(old_cout);
    BOOST_CHECK_EQUAL(ss.str().find("[apple:banana]"), 0);
    BOOST_CHECK_EQUAL(ss.str().find("[orange:grape]"), 0);
}