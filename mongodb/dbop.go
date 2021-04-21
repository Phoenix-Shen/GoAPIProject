package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//连接数据库的语句
var Uri string = "mongodb+srv://root:asdf@ssk.3hxej.mongodb.net/GOAPIPROJDB?retryWrites=true&w=majority&authSource=admin"

//连接数据库操作
// uri:连接语句
// DBName:操作的数据库名称
// timeout:设置超时时间
// numOfConnection:最大链接数量
// collectionName:集合的名字
// 返回:*mongo.Collection操作的集合对象
func ConnectToMongoDB(uri string, DBName string, timeout time.Duration, numOfConnection uint64, collectionName string) (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	o := options.Client().ApplyURI(uri)
	o.SetMaxPoolSize(numOfConnection)
	client, err := mongo.Connect(ctx, o)
	if err != nil {
		return nil, err
	}

	return client.Database(DBName).Collection(collectionName), nil
}

// @title    Create
// @description   插入操作
// @auth      Laozhu
// @param    collection        *mongo.Collection         "操作的Collection"
// @param    obj               []interface{}             "要插入的东西，可以是一个，可以是多个"
// @return    InsertedIDs      []interface{}             "插入的对象的ID，不为nil表明插入成功"
func Create(collection *mongo.Collection, obj []interface{}) []interface{} {
	insertResult, err := collection.InsertMany(context.TODO(), obj)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return insertResult.InsertedIDs
}

// @title    Read
// @description   查找操作
// @auth      Laozhu
// @param    collection        *mongo.Collection         "操作的Collection"
// @param    fliter              interface{}             "筛选器，筛选出要更新的文档"
// @param    outputResult         interface{}            "解码成结构体的输出，需要传引用"
// @return    results      []interface{}                 "查询结果返回值为空或者nil代表没查到和错误，查询返回键值对形式"
func Read(collection *mongo.Collection, fliter interface{}) []interface{} {

	cur, err := collection.Find(context.TODO(), fliter)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var results []interface{}
	for cur.Next(context.TODO()) {
		// 创建一个值，将单个文档解码为该值
		var elem interface{}
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	// 完成后关闭游标
	cur.Close(context.TODO())
	return results
}

// @title    Update
// @description   更新操作
// @auth      Laozhu
// @param    collection        *mongo.Collection         "操作的Collection"
// @param    fliter            interface{}               "筛选器，筛选出要更新的文档"
// @param    obj               interface{}               "更新规则，怎么操作更新"
// @return   updateResult      *mongo.UpdateResult       "更新的结果，如果失败则返回nil"
func Update(collection *mongo.Collection, fliter, obj interface{}) *mongo.UpdateResult {

	//log.Fatal("not implemented Error")
	updateResult, err := collection.UpdateMany(context.TODO(), fliter, obj)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return updateResult
}

// @title    Delete
// @description   删除操作
// @auth      Laozhu
// @param    collection        *mongo.Collection         "操作的Collection"
// @param    fliter            *interface{}              "筛选器，筛选出要更新的文档"
// @return   deleteResult      *mongo.DelteResult        "删除的结果，如果失败则返回nil"
func Delete(collection *mongo.Collection, fliter interface{}) *mongo.DeleteResult {
	deleteResult, err := collection.DeleteMany(context.TODO(), fliter)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return deleteResult
}

//删库跑路
func S删库跑路(collection *mongo.Collection) {
	collection.DeleteMany(context.TODO(), bson.D{{}})
}

/*测试用的代码罢了
func main() {
	collection, _ := connectToMongoDB(uri, "GOAPIPROJDB", 10*time.Second, 5, "User")
	删库跑路(collection)
	result := Read(collection, bson.D{{}})
	print(result)
}
*/
