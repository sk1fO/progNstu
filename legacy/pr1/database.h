#pragma once

#include<iostream>
#include<filesystem>
#include<fstream>
#include"nlohmann/json.hpp"

using namespace std;


string path_create(string scheme_name, string table_name){
    return scheme_name + "/" + table_name;
}

/*string write_values(list& _values){
    string values;
    for(int i; i < _values.size(); i++){
        values += _values[i] -> val + ",";
    }
    values.pop_back();
    return values;
}
*/


void table_init(string& scheme_name, int& tuples_limit){
    nlohmann::json objJson;

    fstream fileInput;
    fileInput.open("scheme.json");
    fileInput >> objJson;
    scheme_name = objJson["name"];
    tuples_limit = objJson["tuples_limit"];

    string table_name;
    for (auto i : objJson["structure"]){
        table_name = i;
        filesystem::create_directories(scheme_name + "/" + table_name);
    }
    fileInput.close();
}


int find_pk_max(string path_to_file){
    string pk_find;
    int pk;
    ifstream readFile(path_to_file);

    if (!getline(readFile, pk_find)){
        return 0;
    }
        


    while(getline(readFile, pk_find)){
        int pkLocation = pk_find.find(',');
        pk_find = pk_find.substr(0, pkLocation);
    }
    readFile.close();
    
    if (pk_find == "pk")
        return 0;

    pk = stoi(pk_find);
    return pk;
}

string path_to_max(string path){

    string maxPath = "1";

    for (const auto& entry : filesystem::directory_iterator(path)){
        string p = entry.path();

        
        if(!(p.find(".csv") == string::npos)){

            int location = p.find('/', path.length());
            p = p.substr(location + 1, p.length());

            int location_dot = p.find(".csv");
            p = p.substr(0, location_dot);
            
            if(p > maxPath)
                maxPath = p;

        }
    }

    return maxPath;
}

bool INSERT_INTO(string _val, string table, string db_name, int tuples_limit){
    string path = path_create(db_name, table);
    string maxPath = path_to_max(path);
    string path_to_file = path + '/' +  maxPath + ".csv";
    int pk = find_pk_max(path_to_file);

    cout << pk << endl;



    if (pk == tuples_limit * stoi(maxPath)){
        path_to_file = path + '/' +  to_string(stoi(maxPath) + 1) + ".csv";
    }

    pk += 1;
    ofstream writeFile;
    writeFile.open(path_to_file,ios_base::app);
    ifstream readFile(path_to_file);
    string line;
    if (!getline(readFile, line)){

        nlohmann::json Jsonread;
        fstream fileInput;
        fileInput.open("scheme.json");
        fileInput >> Jsonread;
        
        string table_header;
        for (string i : Jsonread[table]){
            table_header += i;
            table_header += ",";
        }
        table_header.pop_back();
        fileInput.close();

        writeFile << table_header;
    }

    readFile.close();


    writeFile << "\n" << pk <<"," << _val;
    writeFile.close();

    return 1;

}


bool DELETE_FROM(string path, int _col, string _val){

    for (const auto& entry : filesystem::directory_iterator(path)){
        string p = entry.path();

        
        if(!(p.find(".csv") == string::npos)){
            ifstream readFile;
            readFile.open(p);

            ofstream writeFile;
            writeFile.open("temp.csv");

            
            string _value;
            getline(readFile, _value);
            writeFile << _value;
            while(getline(readFile, _value)){

                string pk, first, second, third, fourth, value;
                value = _value;
                int location;

                location = value.find(',');
                pk = value.substr(0, location);
                value = value.substr(location + 1, value.length());

                location = value.find(',');
                first = value.substr(0, location);
                value = value.substr(location + 1, value.length());

                location = value.find(',');
                second = value.substr(0, location);
                value = value.substr(location + 1, value.length());

                location = value.find(',');
                third = value.substr(0, location);
                value = value.substr(location + 1, value.length());


                fourth = value;


                if(_col == 1){
                    if (_val != first)
                        writeFile << "\n" + _value;
                    
                }

                if(_col == 2 ){
                    if (second != _val)
                        writeFile << "\n" + _value;
                }

                if(_col == 3 ){
                    if(third != _val)
                        writeFile << "\n" + _value;
                }

                if(_col == 4 ){
                    if(fourth != _val)
                        writeFile << "\n" + _value;
                }
                    
            }
            
            readFile.close();
            writeFile.close();
            
            filesystem::remove(p);
            filesystem::copy_file("temp.csv",p);
            filesystem::remove("temp.csv");
        }
        
    }
    return 1;
}

