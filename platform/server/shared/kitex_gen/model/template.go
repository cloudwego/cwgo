// Code generated by thriftgo (0.2.12). DO NOT EDIT.

package model

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"strings"
)

type Template struct {
	Id         int64  `thrift:"id,1" frugal:"1,default,i64" json:"id"`
	Name       string `thrift:"name,2" frugal:"2,default,string" json:"name"`
	Type       int32  `thrift:"type,3" frugal:"3,default,i32" json:"type"`
	IsDeleted  bool   `thrift:"is_deleted,4" frugal:"4,default,bool" json:"is_deleted"`
	CreateTime string `thrift:"create_time,5" frugal:"5,default,string" json:"create_time"`
	UpdateTime string `thrift:"update_time,6" frugal:"6,default,string" json:"update_time"`
}

func NewTemplate() *Template {
	return &Template{}
}

func (p *Template) InitDefault() {
	*p = Template{}
}

func (p *Template) GetId() (v int64) {
	return p.Id
}

func (p *Template) GetName() (v string) {
	return p.Name
}

func (p *Template) GetType() (v int32) {
	return p.Type
}

func (p *Template) GetIsDeleted() (v bool) {
	return p.IsDeleted
}

func (p *Template) GetCreateTime() (v string) {
	return p.CreateTime
}

func (p *Template) GetUpdateTime() (v string) {
	return p.UpdateTime
}
func (p *Template) SetId(val int64) {
	p.Id = val
}
func (p *Template) SetName(val string) {
	p.Name = val
}
func (p *Template) SetType(val int32) {
	p.Type = val
}
func (p *Template) SetIsDeleted(val bool) {
	p.IsDeleted = val
}
func (p *Template) SetCreateTime(val string) {
	p.CreateTime = val
}
func (p *Template) SetUpdateTime(val string) {
	p.UpdateTime = val
}

var fieldIDToName_Template = map[int16]string{
	1: "id",
	2: "name",
	3: "type",
	4: "is_deleted",
	5: "create_time",
	6: "update_time",
}

func (p *Template) Read(iprot thrift.TProtocol) (err error) {

	var fieldTypeId thrift.TType
	var fieldId int16

	if _, err = iprot.ReadStructBegin(); err != nil {
		goto ReadStructBeginError
	}

	for {
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin()
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		switch fieldId {
		case 1:
			if fieldTypeId == thrift.I64 {
				if err = p.ReadField1(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		case 2:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField2(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		case 3:
			if fieldTypeId == thrift.I32 {
				if err = p.ReadField3(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		case 4:
			if fieldTypeId == thrift.BOOL {
				if err = p.ReadField4(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		case 5:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField5(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		case 6:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField6(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		default:
			if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		}

		if err = iprot.ReadFieldEnd(); err != nil {
			goto ReadFieldEndError
		}
	}
	if err = iprot.ReadStructEnd(); err != nil {
		goto ReadStructEndError
	}

	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read struct begin error: ", p), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d begin error: ", p, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d '%s' error: ", p, fieldId, fieldIDToName_Template[fieldId]), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%T field %d skip type %d error: ", p, fieldId, fieldTypeId), err)

ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T read field end error", p), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
}

func (p *Template) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return err
	} else {
		p.Id = v
	}
	return nil
}

func (p *Template) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		p.Name = v
	}
	return nil
}

func (p *Template) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return err
	} else {
		p.Type = v
	}
	return nil
}

func (p *Template) ReadField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return err
	} else {
		p.IsDeleted = v
	}
	return nil
}

func (p *Template) ReadField5(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		p.CreateTime = v
	}
	return nil
}

func (p *Template) ReadField6(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		p.UpdateTime = v
	}
	return nil
}

func (p *Template) Write(oprot thrift.TProtocol) (err error) {
	var fieldId int16
	if err = oprot.WriteStructBegin("Template"); err != nil {
		goto WriteStructBeginError
	}
	if p != nil {
		if err = p.writeField1(oprot); err != nil {
			fieldId = 1
			goto WriteFieldError
		}
		if err = p.writeField2(oprot); err != nil {
			fieldId = 2
			goto WriteFieldError
		}
		if err = p.writeField3(oprot); err != nil {
			fieldId = 3
			goto WriteFieldError
		}
		if err = p.writeField4(oprot); err != nil {
			fieldId = 4
			goto WriteFieldError
		}
		if err = p.writeField5(oprot); err != nil {
			fieldId = 5
			goto WriteFieldError
		}
		if err = p.writeField6(oprot); err != nil {
			fieldId = 6
			goto WriteFieldError
		}

	}
	if err = oprot.WriteFieldStop(); err != nil {
		goto WriteFieldStopError
	}
	if err = oprot.WriteStructEnd(); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", p), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
}

func (p *Template) writeField1(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("id", thrift.I64, 1); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteI64(p.Id); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 end error: ", p), err)
}

func (p *Template) writeField2(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("name", thrift.STRING, 2); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.Name); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 end error: ", p), err)
}

func (p *Template) writeField3(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("type", thrift.I32, 3); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteI32(p.Type); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 3 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 3 end error: ", p), err)
}

func (p *Template) writeField4(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("is_deleted", thrift.BOOL, 4); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteBool(p.IsDeleted); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 4 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 4 end error: ", p), err)
}

func (p *Template) writeField5(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("create_time", thrift.STRING, 5); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.CreateTime); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 5 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 5 end error: ", p), err)
}

func (p *Template) writeField6(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("update_time", thrift.STRING, 6); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.UpdateTime); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 6 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 6 end error: ", p), err)
}

func (p *Template) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Template(%+v)", *p)
}

func (p *Template) DeepEqual(ano *Template) bool {
	if p == ano {
		return true
	} else if p == nil || ano == nil {
		return false
	}
	if !p.Field1DeepEqual(ano.Id) {
		return false
	}
	if !p.Field2DeepEqual(ano.Name) {
		return false
	}
	if !p.Field3DeepEqual(ano.Type) {
		return false
	}
	if !p.Field4DeepEqual(ano.IsDeleted) {
		return false
	}
	if !p.Field5DeepEqual(ano.CreateTime) {
		return false
	}
	if !p.Field6DeepEqual(ano.UpdateTime) {
		return false
	}
	return true
}

func (p *Template) Field1DeepEqual(src int64) bool {

	if p.Id != src {
		return false
	}
	return true
}
func (p *Template) Field2DeepEqual(src string) bool {

	if strings.Compare(p.Name, src) != 0 {
		return false
	}
	return true
}
func (p *Template) Field3DeepEqual(src int32) bool {

	if p.Type != src {
		return false
	}
	return true
}
func (p *Template) Field4DeepEqual(src bool) bool {

	if p.IsDeleted != src {
		return false
	}
	return true
}
func (p *Template) Field5DeepEqual(src string) bool {

	if strings.Compare(p.CreateTime, src) != 0 {
		return false
	}
	return true
}
func (p *Template) Field6DeepEqual(src string) bool {

	if strings.Compare(p.UpdateTime, src) != 0 {
		return false
	}
	return true
}

type TemplateItem struct {
	Id         int64  `thrift:"id,1" frugal:"1,default,i64" json:"id"`
	TemplateId int64  `thrift:"template_id,2" frugal:"2,default,i64" json:"template_id"`
	Name       string `thrift:"name,3" frugal:"3,default,string" json:"name"`
	Content    string `thrift:"content,4" frugal:"4,default,string" json:"content"`
	IsDeleted  bool   `thrift:"is_deleted,5" frugal:"5,default,bool" json:"is_deleted"`
	CreateTime string `thrift:"create_time,6" frugal:"6,default,string" json:"create_time"`
	UpdateTime string `thrift:"update_time,7" frugal:"7,default,string" json:"update_time"`
}

func NewTemplateItem() *TemplateItem {
	return &TemplateItem{}
}

func (p *TemplateItem) InitDefault() {
	*p = TemplateItem{}
}

func (p *TemplateItem) GetId() (v int64) {
	return p.Id
}

func (p *TemplateItem) GetTemplateId() (v int64) {
	return p.TemplateId
}

func (p *TemplateItem) GetName() (v string) {
	return p.Name
}

func (p *TemplateItem) GetContent() (v string) {
	return p.Content
}

func (p *TemplateItem) GetIsDeleted() (v bool) {
	return p.IsDeleted
}

func (p *TemplateItem) GetCreateTime() (v string) {
	return p.CreateTime
}

func (p *TemplateItem) GetUpdateTime() (v string) {
	return p.UpdateTime
}
func (p *TemplateItem) SetId(val int64) {
	p.Id = val
}
func (p *TemplateItem) SetTemplateId(val int64) {
	p.TemplateId = val
}
func (p *TemplateItem) SetName(val string) {
	p.Name = val
}
func (p *TemplateItem) SetContent(val string) {
	p.Content = val
}
func (p *TemplateItem) SetIsDeleted(val bool) {
	p.IsDeleted = val
}
func (p *TemplateItem) SetCreateTime(val string) {
	p.CreateTime = val
}
func (p *TemplateItem) SetUpdateTime(val string) {
	p.UpdateTime = val
}

var fieldIDToName_TemplateItem = map[int16]string{
	1: "id",
	2: "template_id",
	3: "name",
	4: "content",
	5: "is_deleted",
	6: "create_time",
	7: "update_time",
}

func (p *TemplateItem) Read(iprot thrift.TProtocol) (err error) {

	var fieldTypeId thrift.TType
	var fieldId int16

	if _, err = iprot.ReadStructBegin(); err != nil {
		goto ReadStructBeginError
	}

	for {
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin()
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		switch fieldId {
		case 1:
			if fieldTypeId == thrift.I64 {
				if err = p.ReadField1(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		case 2:
			if fieldTypeId == thrift.I64 {
				if err = p.ReadField2(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		case 3:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField3(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		case 4:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField4(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		case 5:
			if fieldTypeId == thrift.BOOL {
				if err = p.ReadField5(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		case 6:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField6(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		case 7:
			if fieldTypeId == thrift.STRING {
				if err = p.ReadField7(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		default:
			if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		}

		if err = iprot.ReadFieldEnd(); err != nil {
			goto ReadFieldEndError
		}
	}
	if err = iprot.ReadStructEnd(); err != nil {
		goto ReadStructEndError
	}

	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read struct begin error: ", p), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d begin error: ", p, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d '%s' error: ", p, fieldId, fieldIDToName_TemplateItem[fieldId]), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%T field %d skip type %d error: ", p, fieldId, fieldTypeId), err)

ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T read field end error", p), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
}

func (p *TemplateItem) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return err
	} else {
		p.Id = v
	}
	return nil
}

func (p *TemplateItem) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return err
	} else {
		p.TemplateId = v
	}
	return nil
}

func (p *TemplateItem) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		p.Name = v
	}
	return nil
}

func (p *TemplateItem) ReadField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		p.Content = v
	}
	return nil
}

func (p *TemplateItem) ReadField5(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return err
	} else {
		p.IsDeleted = v
	}
	return nil
}

func (p *TemplateItem) ReadField6(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		p.CreateTime = v
	}
	return nil
}

func (p *TemplateItem) ReadField7(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return err
	} else {
		p.UpdateTime = v
	}
	return nil
}

func (p *TemplateItem) Write(oprot thrift.TProtocol) (err error) {
	var fieldId int16
	if err = oprot.WriteStructBegin("TemplateItem"); err != nil {
		goto WriteStructBeginError
	}
	if p != nil {
		if err = p.writeField1(oprot); err != nil {
			fieldId = 1
			goto WriteFieldError
		}
		if err = p.writeField2(oprot); err != nil {
			fieldId = 2
			goto WriteFieldError
		}
		if err = p.writeField3(oprot); err != nil {
			fieldId = 3
			goto WriteFieldError
		}
		if err = p.writeField4(oprot); err != nil {
			fieldId = 4
			goto WriteFieldError
		}
		if err = p.writeField5(oprot); err != nil {
			fieldId = 5
			goto WriteFieldError
		}
		if err = p.writeField6(oprot); err != nil {
			fieldId = 6
			goto WriteFieldError
		}
		if err = p.writeField7(oprot); err != nil {
			fieldId = 7
			goto WriteFieldError
		}

	}
	if err = oprot.WriteFieldStop(); err != nil {
		goto WriteFieldStopError
	}
	if err = oprot.WriteStructEnd(); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", p), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
}

func (p *TemplateItem) writeField1(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("id", thrift.I64, 1); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteI64(p.Id); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 end error: ", p), err)
}

func (p *TemplateItem) writeField2(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("template_id", thrift.I64, 2); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteI64(p.TemplateId); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 end error: ", p), err)
}

func (p *TemplateItem) writeField3(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("name", thrift.STRING, 3); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.Name); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 3 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 3 end error: ", p), err)
}

func (p *TemplateItem) writeField4(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("content", thrift.STRING, 4); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.Content); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 4 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 4 end error: ", p), err)
}

func (p *TemplateItem) writeField5(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("is_deleted", thrift.BOOL, 5); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteBool(p.IsDeleted); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 5 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 5 end error: ", p), err)
}

func (p *TemplateItem) writeField6(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("create_time", thrift.STRING, 6); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.CreateTime); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 6 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 6 end error: ", p), err)
}

func (p *TemplateItem) writeField7(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("update_time", thrift.STRING, 7); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteString(p.UpdateTime); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 7 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 7 end error: ", p), err)
}

func (p *TemplateItem) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TemplateItem(%+v)", *p)
}

func (p *TemplateItem) DeepEqual(ano *TemplateItem) bool {
	if p == ano {
		return true
	} else if p == nil || ano == nil {
		return false
	}
	if !p.Field1DeepEqual(ano.Id) {
		return false
	}
	if !p.Field2DeepEqual(ano.TemplateId) {
		return false
	}
	if !p.Field3DeepEqual(ano.Name) {
		return false
	}
	if !p.Field4DeepEqual(ano.Content) {
		return false
	}
	if !p.Field5DeepEqual(ano.IsDeleted) {
		return false
	}
	if !p.Field6DeepEqual(ano.CreateTime) {
		return false
	}
	if !p.Field7DeepEqual(ano.UpdateTime) {
		return false
	}
	return true
}

func (p *TemplateItem) Field1DeepEqual(src int64) bool {

	if p.Id != src {
		return false
	}
	return true
}
func (p *TemplateItem) Field2DeepEqual(src int64) bool {

	if p.TemplateId != src {
		return false
	}
	return true
}
func (p *TemplateItem) Field3DeepEqual(src string) bool {

	if strings.Compare(p.Name, src) != 0 {
		return false
	}
	return true
}
func (p *TemplateItem) Field4DeepEqual(src string) bool {

	if strings.Compare(p.Content, src) != 0 {
		return false
	}
	return true
}
func (p *TemplateItem) Field5DeepEqual(src bool) bool {

	if p.IsDeleted != src {
		return false
	}
	return true
}
func (p *TemplateItem) Field6DeepEqual(src string) bool {

	if strings.Compare(p.CreateTime, src) != 0 {
		return false
	}
	return true
}
func (p *TemplateItem) Field7DeepEqual(src string) bool {

	if strings.Compare(p.UpdateTime, src) != 0 {
		return false
	}
	return true
}

type TemplateWithInfo struct {
	Template *Template       `thrift:"template,1" frugal:"1,default,Template" json:"template"`
	Items    []*TemplateItem `thrift:"items,2" frugal:"2,default,list<TemplateItem>" json:"items"`
}

func NewTemplateWithInfo() *TemplateWithInfo {
	return &TemplateWithInfo{}
}

func (p *TemplateWithInfo) InitDefault() {
	*p = TemplateWithInfo{}
}

var TemplateWithInfo_Template_DEFAULT *Template

func (p *TemplateWithInfo) GetTemplate() (v *Template) {
	if !p.IsSetTemplate() {
		return TemplateWithInfo_Template_DEFAULT
	}
	return p.Template
}

func (p *TemplateWithInfo) GetItems() (v []*TemplateItem) {
	return p.Items
}
func (p *TemplateWithInfo) SetTemplate(val *Template) {
	p.Template = val
}
func (p *TemplateWithInfo) SetItems(val []*TemplateItem) {
	p.Items = val
}

var fieldIDToName_TemplateWithInfo = map[int16]string{
	1: "template",
	2: "items",
}

func (p *TemplateWithInfo) IsSetTemplate() bool {
	return p.Template != nil
}

func (p *TemplateWithInfo) Read(iprot thrift.TProtocol) (err error) {

	var fieldTypeId thrift.TType
	var fieldId int16

	if _, err = iprot.ReadStructBegin(); err != nil {
		goto ReadStructBeginError
	}

	for {
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin()
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRUCT {
				if err = p.ReadField1(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		case 2:
			if fieldTypeId == thrift.LIST {
				if err = p.ReadField2(iprot); err != nil {
					goto ReadFieldError
				}
			} else {
				if err = iprot.Skip(fieldTypeId); err != nil {
					goto SkipFieldError
				}
			}
		default:
			if err = iprot.Skip(fieldTypeId); err != nil {
				goto SkipFieldError
			}
		}

		if err = iprot.ReadFieldEnd(); err != nil {
			goto ReadFieldEndError
		}
	}
	if err = iprot.ReadStructEnd(); err != nil {
		goto ReadStructEndError
	}

	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read struct begin error: ", p), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d begin error: ", p, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%T read field %d '%s' error: ", p, fieldId, fieldIDToName_TemplateWithInfo[fieldId]), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%T field %d skip type %d error: ", p, fieldId, fieldTypeId), err)

ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T read field end error", p), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
}

func (p *TemplateWithInfo) ReadField1(iprot thrift.TProtocol) error {
	p.Template = NewTemplate()
	if err := p.Template.Read(iprot); err != nil {
		return err
	}
	return nil
}

func (p *TemplateWithInfo) ReadField2(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return err
	}
	p.Items = make([]*TemplateItem, 0, size)
	for i := 0; i < size; i++ {
		_elem := NewTemplateItem()
		if err := _elem.Read(iprot); err != nil {
			return err
		}

		p.Items = append(p.Items, _elem)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return err
	}
	return nil
}

func (p *TemplateWithInfo) Write(oprot thrift.TProtocol) (err error) {
	var fieldId int16
	if err = oprot.WriteStructBegin("TemplateWithInfo"); err != nil {
		goto WriteStructBeginError
	}
	if p != nil {
		if err = p.writeField1(oprot); err != nil {
			fieldId = 1
			goto WriteFieldError
		}
		if err = p.writeField2(oprot); err != nil {
			fieldId = 2
			goto WriteFieldError
		}

	}
	if err = oprot.WriteFieldStop(); err != nil {
		goto WriteFieldStopError
	}
	if err = oprot.WriteStructEnd(); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", p), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
}

func (p *TemplateWithInfo) writeField1(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("template", thrift.STRUCT, 1); err != nil {
		goto WriteFieldBeginError
	}
	if err := p.Template.Write(oprot); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 1 end error: ", p), err)
}

func (p *TemplateWithInfo) writeField2(oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteFieldBegin("items", thrift.LIST, 2); err != nil {
		goto WriteFieldBeginError
	}
	if err := oprot.WriteListBegin(thrift.STRUCT, len(p.Items)); err != nil {
		return err
	}
	for _, v := range p.Items {
		if err := v.Write(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return err
	}
	if err = oprot.WriteFieldEnd(); err != nil {
		goto WriteFieldEndError
	}
	return nil
WriteFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 begin error: ", p), err)
WriteFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%T write field 2 end error: ", p), err)
}

func (p *TemplateWithInfo) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TemplateWithInfo(%+v)", *p)
}

func (p *TemplateWithInfo) DeepEqual(ano *TemplateWithInfo) bool {
	if p == ano {
		return true
	} else if p == nil || ano == nil {
		return false
	}
	if !p.Field1DeepEqual(ano.Template) {
		return false
	}
	if !p.Field2DeepEqual(ano.Items) {
		return false
	}
	return true
}

func (p *TemplateWithInfo) Field1DeepEqual(src *Template) bool {

	if !p.Template.DeepEqual(src) {
		return false
	}
	return true
}
func (p *TemplateWithInfo) Field2DeepEqual(src []*TemplateItem) bool {

	if len(p.Items) != len(src) {
		return false
	}
	for i, v := range p.Items {
		_src := src[i]
		if !v.DeepEqual(_src) {
			return false
		}
	}
	return true
}
