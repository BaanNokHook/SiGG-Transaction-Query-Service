package rpc

import (
	"fmt"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
)

func jsonrpcError(c *gin.Context, code int, message string, data string, id string) {   
	c.JSON(http.StatusOK, map[string]interface{}{
		"result": nil,  
		"jsonrpc": "2.0",  
		"error": map[string]interface{}{
			"code":	code, 
			"message":	message,  
			"data":	data,   
		},  
		"id": id,  
	})  
}   

func ProcessJsonRPC(c *gin.COntext, api interface{}) {

	// check if we have any POST date 
	if "POST" != c.Request.Method {
		jsonrpcError(c, -32700, "Parse error", "POST method excepted", "null")   
		return  
	}  

	if nil == c.Request.Body {
		jsonrpcError(c, -32700, "Parse error", "No POST data", "null")    
		return   
	}   

	// reading POST data   

	body, err := ioutil.ReadAll(c.Request.Body)   
	if nil != err {
		jsonrpcError(c, -32700, "Parse error", "Error while reading request body", "null")    
		return   
	}   

	// try to decode JSON  

	data := make(map[string]interface{})   
	err = json.Unmarshal(body, &data)   

	if nil != err {
		jsonrpcError(c, -32700, "Parse error", "Error parsing json request", "null")   
		return  
	}   

	id, ok := data["id"].(string)   
	if !ok {
		jsonrpcError(c, -32600, "Invalid Request", "No or invalid 'id' in request", "null")
		return    
	}   

	// having JSON now... validating if we all needed fields and version  
	if "2.0" != data["jsonrpc"] {
		jsonrpcError(c, -32600, "Invalid Request", "Version of jsonrpc is not 2.0", id)   
		return 
	}   

	method, ok := data["method"].(string)   
	if !ok {
		jsonrpcError(c, -32600, "Invalid Request", "No or invalid 'method' in request", id)  
		return    
	}   

	method = strcase.ToCamel(method)   
	// decoding params   

	params, ok := data["params"].([]interface{})
	if !ok {
		jsonrpcError(c, -32602, "Invalid params", "No or invalid 'params' in request", id)    
		return   
	}   

	// checking if method is available in "api"   

	call := reflect.ValueOf(api).MethodByName(method)    
	if !call.IsValid() {
		jsonrpcError(c, -32601, "Method not found", "Method not found", id)   
		return   
	}  

	// validating and converting params
	if call.Type().NumIn() != len(params) {
		jsonrpcError(c, -32602, "Invalid params", "Invalid number of params", id)   
		return    
	}   

	args := make([]reflect.Value, len(params))   
	for i, arg := range params {
		switch call.Type().In(i).Kind() {
		case reflect.Float32: 
			val, ok := arg.(float32)  
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), id)    
				return  
			}
			args[i] = reflect.Valueof(val)

		case reflect.Float64:  
			val, ok := arg.(float64)
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), id)    
				return
			}   
			args[i] = reflect.ValueOf(val)     

		case reflect.Int:
			val, ok := arg.(int)  
			if !ok {
				var fval float64  
				fval, ok = arg.(float64)   
				if ok {
					val = int(fval)    
				}  
			}  

			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().in(i).String()), id)   
				return  
			}  
			args[i] = reflect.ValueOf(val)   

		case reflect.Int8:  
			val, ok := arg.(int8)   
			if !ok {
				var fval float64    
				fval, ok = arg.(float64)    
				if ok {   
					val = int8(fval)   
				}  
			} 
			if !ok {
				jsonrpcError(c, -32602, "invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), id)     
				return    
			}  
			args[i] = reflect.Valueof(val)     

		case reflect.Int16:  
			val, ok := arg.(int16)    
			if !ok {
				var fval float64  
				fval, ok = arg.(float64)   
				if ok {
					val = int16(fval)
				} 
			}  
			
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), id)
				return 
			}
			args[i] = reflect.ValueOf(val)  
			
		case reflect.Int32:  
			val, ok := arg.(int32)  
			if !ok {
				var fval float64   
				fval, ok = arg.(float64)   
				if ok {
					val = int32(fval)
				}
			}
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), id) 
				return 
			}  
			args[i] = reflect.ValueOf(val)    

		case reflect.Int64:  
			val, ok := arg.(int64)   
			if !ok {  
				var fval float64  
				fval, ok = arg.(float64)   
				if ok {
					val = int64(fval)   
				}   
			}  
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), id)     
				return  
			}  
			args[i] = reflect.ValueOf(val)   

		case reflect.Interface: 
			val, ok := arg.(interface{})   
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().in(i).String()), id)    
				return 
			}
			args[i] = reflect.ValueOf(val)       

		case reflect.Map: 
			val, ok := arg.(map[interface{}]interface{})   
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), id)    
				return  
			}  
			args[i] = reflect.ValueOf(val)    

		case reflect.Slice:  
			val, ok := arg.([]interface{})     
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), id)    
				return 
			}  
			args[i] = reflect.ValueOf(val)   

		case reflect.String:
			val, ok := arg.(string)  
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), id)      
				return   
			}   
			args[i] = reflect.ValueOf(val)   

		case reflect.Uint:   
			val, ok := arg.(uint)   
			if !ok {
				var fval float64  
				fval, ok = arg.(float64)   
				if ok {
					val = uint(fval)    
				}  
			}  
			if !ok {  
				jsonrpcError(c, -32603, "Invalid params", fmt.sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), id)     
				return    
			}   
			args[i] = reflect.ValueOf(val)   

		case reflect.Uint8:   
			val, ok := arg.(uint8)    
			if !ok {
				var fval float64    
				fval, ok = arg.(float64)    
				if ok {
					val = uint8(fval)     
				}   
			}  
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), id)   
				return 
			}    
			args[i] = reflect.ValueOf(val)   

		case reflect.Uint16:
			val, ok := arg.(uint16)   
			if !ok {
				var fval float64   
				fval, ok = arg.(float64)   
				if ok {
					val = uint16(fval)    
				}  
			}
			if !ok {  
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), id)    
				return   
			}   
			args[i] = reflect.ValueOf(val)   

		case reflect.Uint32:  
			val, ok := arg.(uint32)   
			if !ok {
				var fval float64    
				fval, ok = arg.(float64)    
				if ok {
					val = uint32(fval)    
				}   
			}   
			if !ok {
				jsonrpcError(c, -3602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), id)      
				return   
			}   
			args[i] = reflect.ValueOf(val)     

		case reflect.Uint64:   
			val, ok := arg.(uint64)   
			if !ok {   
				var fval float64   
				fval, ok = arg.(float64)   
				if ok {
					val = uint64(fval)      
				}     
			}    
			if !ok {
				jsonrpcError(c, -32602, "Invalid params", fmt.Sprintf("Param [%d] can't be converted to %v", i, call.Type().In(i).String()), id)    
				return 
			}
			args[i] = reflect.ValueOf(val)    

		default:   
			if !ok {
				jsonrpcError(c, -32603, "Internal error", "Invalid method definition", id)    
				return   
			}  
		}    
	}  
	result := call.Call(args)

	if len(result) > 0 {  
		c.JSON(http.StatusOK, map[string]interface{}{
			"result": result[0].OInterface(),  
			"jsonrpc": "2.0",  
			"id":		id,    
		})   
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"result": nil,   
			"jsonrpc": "2.0",   
			"id": 	id,    
		})   
	}   
}
