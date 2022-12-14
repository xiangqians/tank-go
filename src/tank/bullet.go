// bullet
// 子弹
//
// @author xiangqian
// @date 11:01 2022/11/27
package tank

import (
	"github.com/hajimehoshi/ebiten"
	"time"
)

type Bullet struct {
	*AbsGraphics
	TankId string // 子弹所属坦克id
}

func CreateBullet(pTank *Tank, speed Speed) *Bullet {
	// 让子弹坐标从坦克中心发出
	pLocation := pTank.Location
	width, height := pTank.pImage.Size()
	//log.Printf("width = %v, height = %v\n", width, height)
	var location Location
	switch pTank.Direction {
	case DirectionUp:
		location = Location{pLocation.X + float64(width/2) - 4.51392, pLocation.Y - 8}
	case DirectionDown:
		location = Location{pLocation.X + float64(width/2) - 5.8, pLocation.Y + float64(height)}
	case DirectionLeft:
		location = Location{pLocation.X - 8, pLocation.Y + float64(height/2) - 2.91392}
	case DirectionRight:
		location = Location{pLocation.X + float64(width), pLocation.Y + float64(height/2) - 2.91392}
	}
	//log.Printf("location = %v\n", location)

	pBullet := &Bullet{
		AbsGraphics: CreateAbsGraphics("", GraphicsTyBullet, location, pTank.Direction, speed),
		TankId:      pTank.Id,
	}
	pBullet.Init()
	return pBullet
}

func (pBullet *Bullet) Init() {
	pBullet.AbsGraphics.Init(pBullet)
}

func (pBullet *Bullet) Intersect(x, y float64, otherGraphics Graphics) bool {
	r := pBullet.AbsGraphics.Intersect(x, y, otherGraphics)
	if !r {
		return r
	}

	// 子弹打中子弹，则双方子弹抵消
	if otherGraphics.GetGraphicsTy() == GraphicsTyBullet {
		pBullet.Status = StatusTerm
	} else
	// 子弹打中坦克
	if otherGraphics.GetGraphicsTy() == GraphicsTyTank {
		pTank := otherGraphics.(*Tank)
		if pTank.Status != StatusTerm {
			pTank.Hp--
		}

		if pTank.Hp <= 0 {
			pTank.Status = StatusTerm
		}

		pApp.pEndpoint.SendGraphicsToAddrs(otherGraphics)
	}

	return r
}

func (pBullet *Bullet) IsOutOfBounds(x, y float64) bool {
	r := pBullet.AbsGraphics.IsOutOfBounds(x, y)
	if r {
		pBullet.Status = StatusTerm
		//log.Printf("Bullet OutOfBounds, %v\n", pBullet.GetId())
		pApp.pGame.DelGraphics(pBullet)
	}
	return r
}

func (pBullet *Bullet) Draw(screen *ebiten.Image) error {
	if err := pBullet.AbsGraphics.Draw(screen); err != nil {
		return err
	}

	//if pBullet.TankId == pApp.pGame.pTank.Id && pBullet.Status == StatusNew {
	//}

	return nil
}

// 如果是当前用户所发射的子弹，那么由当前用户轮询设置子弹位置
func (pBullet *Bullet) Run() {
	if pBullet.Status == StatusNew {
		pBullet.Status = StatusRun
	}

	for {
		pBullet.Move(pBullet.Direction)
		pApp.pEndpoint.SendGraphicsToAddrs(pBullet)

		if pBullet.Status != StatusRun {
			break
		}

		// Speed value
		var speedValue int8
		switch pBullet.Speed {
		case SpeedSlow:
			speedValue = 10
		case SpeedNormal:
			speedValue = 5
		case SpeedFast:
			speedValue = 2
		default:
			speedValue = 16
		}

		time.Sleep(time.Duration(speedValue) * time.Millisecond)
	}
	//log.Printf("%v Term\n", pBullet.id)
}

func (pBullet *Bullet) UpImage() *ebiten.Image {
	return pApp.pImage.pBulletUpImage
}

func (pBullet *Bullet) DownImage() *ebiten.Image {
	return pApp.pImage.pBulletDownImage
}

func (pBullet *Bullet) LeftImage() *ebiten.Image {
	return pApp.pImage.pBulletLeftImage
}

func (pBullet *Bullet) RightImage() *ebiten.Image {
	return pApp.pImage.pBulletRightImage
}
