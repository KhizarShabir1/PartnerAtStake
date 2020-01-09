package assignment02IBC
import
(
  "crypto/sha256"
  "encoding/hex"
  "fmt"
  "strconv"
  // "reflect"
  "math/rand"
  "time"
)

type TransSend struct{
  Transactions string
  FreeCoin int  //free for life
  Sender string //name off sender
  Miner string
}

type Trans struct{
  Transactions []string
  FreeCoin []int  //free for life
  NoOfTrans int
}

type Block struct {
Transaction *Trans
PrevPointer *Block
HashVal string
}

func InsertBlock( transaction Trans,chainHead *Block) (string ,*Block) { //returning blockchain and hash value of new added block

  if chainHead==nil{
    block:=Block{&transaction,nil, "emp"}
    chainHead:= & block
    return "Genesis Block",chainHead
  } else{

    block:=Block{&transaction,nil, "emp"}
    newBlock := block
    prev :=chainHead
    curr:=chainHead
    // xt:=reflect.TypeOf(curr).Kind()
    // fmt.Printf("%T: %s\n", xt, xt)

      for curr!=nil{
        prev=curr
        curr=curr.PrevPointer
    	}

    prev.PrevPointer=&newBlock
    hash := sha256.New()
    var tempStr string
    tempStr=""
    for i:=0;i<len(prev.Transaction.Transactions);i++ {

        tempStr+=prev.Transaction.Transactions[i]
    }
    for i:=0;i<len(prev.Transaction.FreeCoin);i++ {

        tempStr+=strconv.Itoa(prev.Transaction.FreeCoin[i])

    }
    rand.Seed(time.Now().UnixNano())
    nk := 100 + rand.Intn(512669-100+1)
    tempStr+=strconv.Itoa(nk)


    hash.Write([]byte(tempStr))
    prev.PrevPointer.HashVal= hex.EncodeToString(hash.Sum(nil))
    return prev.PrevPointer.HashVal,chainHead
  }
}
func ListBlocks(chainHead *Block) {

  curr:=chainHead
  fmt.Printf("\n")
  for curr!=nil{
    for j:=0;j<len(curr.Transaction.Transactions);j++{
      fmt.Printf(curr.Transaction.Transactions[j])
      fmt.Print(curr.Transaction.FreeCoin[j])
      fmt.Printf("  -> ")
    }


    curr=curr.PrevPointer
  }
  fmt.Printf("\n")

}

func CheckHashExists(chainHead *Block, hashVal string)bool {

  curr:=chainHead

  for curr!=nil{

    if curr.HashVal==hashVal{
      fmt.Println("Hash already exists")

      return true
    }

    curr=curr.PrevPointer
  }
  fmt.Println("Hash does not exists")
  return false



}




// func ChangeBlock(oldTrans string, newTrans string, chainHead *Block) {
//   curr:=chainHead
//   for curr!=nil{
//
//     if curr.Transaction==oldTrans{
//       curr.Transaction=newTrans
//     }
//
//     curr=curr.PrevPointer
//   }
//
// }
// func VerifyChain(chainHead *Block) {
//
//   curr:=chainHead
//
//   check:= true
//
//   for curr!=nil{
//
//     hash := sha256.New()
//     hash.Write([]byte(curr.Transaction))
//     hashStr:= hex.EncodeToString(hash.Sum(nil))
//     if curr.PrevPointer!=nil{
//       if hashStr!=curr.PrevPointer.HashVal{
//         check= false
//       }
//     }
//
//
//     curr=curr.PrevPointer
//   }
//   if check{
//
//     fmt.Printf("\nBlockChain verified !!!\n")
//   }else{
//     fmt.Printf("\nBlockChain NOT verified !!!\n")
//   }
//
//
// }
