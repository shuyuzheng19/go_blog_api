package requests

type Sort string

const (
	CREATE Sort = "CREATE" //通过创建日期排序
	UPDATE Sort = "UPDATE" //通过修改日期排序
	BACK   Sort = "BACK"   //通过创建日期倒叙
	ID     Sort = "ID"     //通过ID排序
	EYE    Sort = "EYE"    //浏览量排序
	SIZE   Sort = "SIZE"   //文件大小正序
	BSIZE  Sort = "BSIZE"  //文件大小倒叙
)

type RequestQuery struct {
	Page int  `form:"page"`
	Size int  `form:"size"`
	Cid  *int `form:"cid"`
	Sort Sort `form:"sort"`
	Tid  *int `form:"tid"`
	Uid  *int `form:"uid"`
}
