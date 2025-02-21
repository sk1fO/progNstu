#include <iostream>
#include <ctime>
#include <cmath>
#include <vector>
#include <random>
#include <set>
using namespace std;

int modbin(int a,int x,int m){
    int result = 1;
    int mod = a % m;
    while ( x != 0 ){
        if (x%2==1) result = (result*mod)%m;
            mod = (mod*mod)%m;
            x=x/2;
    }
    return result;
}

int pr_miller(vector<int> razloj){
    int m=1;
    for(int i:razloj)m*=i;
    int n=2*m+1;
    return n;
}

int pr_pokling(vector<int> razloj){
    int F=1;
    for(int i:razloj)F*=i;
    
    int R=F>>1;
    R-=R%2;
    int n=R*F+1;

    return n;
}

int is_pr_miller(int prime,int t,vector<int> razloj){
    razloj.push_back(2);
    srand(time(0));
    vector <int> randoms;
    
    for (int i=0;i<t;i++)randoms.push_back(rand()%(prime-2)+1);

    for (int j:randoms)if (modbin(j,prime-1,prime)!=1) return 0;

    bool is_one_flag;
    for (int q:razloj){
        is_one_flag=true;
        for(int a:randoms){
            if (modbin(a,(prime-1)/q,prime)!=1){
                is_one_flag=false;
                break;
            }
        }
        if (is_one_flag)return -1;
    }

    return 1;
}

int is_pr_pokling(int prime,int t,vector<int> razloj){
    srand(time(0));
    vector <int> randoms;

    for (int i=0;i<t;i++)randoms.push_back(rand()%(prime-2)+1);
    for (int j:randoms)if (modbin(j,prime-1,prime)!=1) return 0;

    bool is_zero_flag;
    for (int a:randoms){
        is_zero_flag=true;

        for(int q:razloj){
            if (modbin(a,(prime-1)/q,prime)==1){
                is_zero_flag=false;
                break;
            }
        }
        if (is_zero_flag)return 1;
    }
    return -1;
}

int GOST(int q,int t){
    int N,u,p,step=1;
    double E;
    auto degree=[](int num,int degr){double res=1;while(degr>0){res*=num,degr--;} return res;};

    while(true){
    switch(step) {

    case 1:    //1
    E=0;
    N=ceil(degree(2,t-1)/q) + ceil(degree(2,t-1)*E/q);
    if(N%2==1)N++;

    u=0;   //2

    case 3:     //3
    p=(N+u)*q+1;

    if(p>degree(2,t)){  //4
        step=1;
        break;
    }

    if(modbin(2,p-1,p)==1 and modbin(2,N+u,p)!=1) return p;

    u+=2;
    step=3;
    }
    }
}



int main(){
int prime,t,k;
vector<int> resheto={2};
bool is_prime;
for (int i=3;i<500;i++){
    is_prime=true;
    for (int j:resheto){
        if(i%j==0){
            is_prime=false;
            break;
        }
    }
    if(is_prime){
        resheto.push_back(i);
    }
}
cout<<"Resheto"<<endl;
for (auto j:resheto) cout << j << " ";
cout<<endl;

vector<int> miller_primes,poklin_primes;
vector<vector<int>> razloj(10);
cout<<"Miller"<<endl;
for(int i=0;i<10;i++)cout.width(6),cout<<i+1;
cout<<endl;
for(int i=0;i<10;i++){
    razloj[i]={resheto[i+1],resheto[i+2]};
    miller_primes.push_back(pr_miller(razloj[i]));
    cout.width(6);
    cout<<miller_primes[i];
}
cout<<endl;
for(int i=0;i<10;i++){
    cout.width(6);
    cout<<is_pr_miller(miller_primes[i],10,razloj[i]);
}
cout<<endl;

cout<<"Poklington"<<endl;
for(int i=0;i<10;i++)cout.width(6),cout<<i+1;
cout<<endl;
for(int i=0;i<10;i++){
    razloj[i]={resheto[i+3]};
    poklin_primes.push_back(pr_pokling(razloj[i]));
    cout.width(6);
    cout<<poklin_primes[i];
}
cout<<endl;
for(int i=0;i<10;i++){
    cout.width(6);
    cout<<is_pr_pokling(poklin_primes[i],10,razloj[i]);
}
cout<<endl;


cout<<"GOST numbers"<<endl;
cout<<GOST(3,4)<<endl;
cout<<GOST(5,6)<<endl;
cout<<GOST(7,5)<<endl;
cout<<GOST(5,5)<<endl;
cout<<GOST(11,7)<<endl;
cout<<GOST(11,8)<<endl;
cout<<GOST(13,7)<<endl;
cout<<GOST(13,8)<<endl;
cout<<GOST(17,9)<<endl;
cout<<GOST(17,10)<<endl;
}