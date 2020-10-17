// Package main Implementación del servidor para el servicio Paquete
package main

import (
	"context"
	"log"
	"net"
	"math/rand"
	"time"
	"strings"
	"strconv"
	"fmt"
	"os"
	"encoding/csv"
	// Llamamos el paquete de gRPC
	"google.golang.org/grpc"

	// Llamamos el compilado que nos generó protoc
	pb "./logistica"
)

const (
	// Puerto por donde expondremos los servicios del cliente
	port = ":50053"
	// Puerto por donde expondremos los servicios del reparto y sus camiones
	port2 = ":50054"
)

// server es usado para implementar Paquete.
type server struct{}

var codigos []string
var n = 0
var colaRetail = make([]string, 0)
var colaNormal = make([]string, 0)
var colaPrioritario = make([]string, 0)


type Paquete struct {
	Id string
	Tipo int32
	Nombre string
	Valor int32
	Origen string
	Destino string
	Seguimiento string
	Timestamp string
	Estado string
	Intentos int32
	Intentar bool//si es true se intentara entregar, sino no es factible pues no se gana nada
}
//estructura que contendra todos las ordenes de envio
var paquetes []Paquete

func tomarPaq(id string) int {//solo se toma el paquete si aun es factible entregarlo (si se gana algo)
	it := 0
	for i := 0; i < len(paquetes); i++ {
		if paquetes[i].Id == id{
			it = i
			break
		}
	}
	return it
}

func entregarPaq(idCamion int32) (int, bool) {
	if idCamion == 1 || idCamion == 2{
		if len(colaRetail) > 0 {
			x := colaRetail[0]
			colaRetail = colaRetail[1:] // Se descarta el top de la cola
			it := tomarPaq(x)
			if paquetes[it].Intentar == true{
				paquetes[it].Estado = "En camino"
				paquetes[it].Intentos += 1 
				return it,true
			}
		}else if len(colaPrioritario) > 0 {
			x := colaPrioritario[0]
			colaPrioritario = colaPrioritario[1:] // Se descarta el top de la cola
			it := tomarPaq(x)
			if paquetes[it].Intentar == true{
				paquetes[it].Estado = "En camino"
				paquetes[it].Intentos += 1 
				return it,true 
			}
		}
	}else if idCamion == 3{
		if len(colaRetail) > 0 {
			x := colaRetail[0]
			colaRetail = colaRetail[1:] // Se descarta el top de la cola
			it := tomarPaq(x)
			if paquetes[it].Intentar == true{
				paquetes[it].Estado = "En camino"
				paquetes[it].Intentos += 1 
				return it,true
			}
		}else if len(colaNormal) > 0 {
			x := colaNormal[0]
			colaNormal = colaNormal[1:] // Se descarta el top de la cola
			it := tomarPaq(x)
			if paquetes[it].Intentar == true{
				paquetes[it].Estado = "En camino"
				paquetes[it].Intentos += 1 
				return it,true 
			}
		}
	}

	//return p
	return 0,false
}

func generarId() string {
	s := ""
	l := []string{"A", "B", "C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V"}
	n := []string{"0","1","2","3","4","5","6","7","8","9"}
	
	for{
		t := true
		var randomI int 
		randomI = rand.Intn(len(l)) 
		s += l[randomI]
		randomI = rand.Intn(len(l))
		s += l[randomI]
		randomI = rand.Intn(len(n))
		s += n[randomI]
		randomI = rand.Intn(len(n))
		s += n[randomI]
		randomI = rand.Intn(len(n))
		s += n[randomI]
		randomI = rand.Intn(len(n))
		s += n[randomI]

		for i := range codigos {
			if s == codigos[i]{
				t=false
			}
		}
		if t {
			codigos = append(codigos, s)
			break
		}
	}
	return s
}

func generarId2(nn int) string{//retorna el id del paquete que sera un string numerico de 5 digitos
	var s = ""
	var sN = strconv.Itoa(nn)
	s = strings.Repeat("0", 5-len(sN))
	s += sN
	return s
}

func (s *server) InformeReparto(ctx context.Context, in *pb.InformeRequest) (*pb.InformeResponse, error) {
	log.Println("Recibiendo informes de entrega desde reparto para el paquete", in.Id)
	for i := 0; i < len(paquetes); i++ {
		if paquetes[i].Id ==  in.Id{
			if in.Entregado{
				paquetes[i].Estado = "Recibido"
				paquetes[i].Intentar = false
			}else{
				if paquetes[i].Tipo == 0{
					if paquetes[i].Intentos < 3{
						//si es factible entregarlo denuevo, se añade a la respectiva cola
						paquetes[i].Estado = "En bodega"
						colaRetail = append(colaRetail, paquetes[i].Id)
					}else{
						paquetes[i].Estado = "No Recibido"
						paquetes[i].Intentar = false
					}

				}else if paquetes[i].Tipo == 1{
					if paquetes[i].Intentos < 2 && paquetes[i].Intentos*10 <= paquetes[i].Valor{
						//si es factible entregarlo denuevo, se añade a la respectiva cola
						paquetes[i].Estado = "En bodega"
						colaNormal = append(colaNormal, paquetes[i].Id)
					}else{
						paquetes[i].Estado = "No Recibido"
						paquetes[i].Intentar = false
					}
				}else if paquetes[i].Tipo == 2{
					if paquetes[i].Intentos < 2 && paquetes[i].Intentos*10 <= paquetes[i].Valor{
						//si es factible entregarlo denuevo, se añade a la respectiva cola
						paquetes[i].Estado = "En bodega"
						colaPrioritario = append(colaPrioritario, paquetes[i].Id)
					}else{
						paquetes[i].Estado = "No Recibido"
						paquetes[i].Intentar = false
					}
				}
			}
			break
		}
	}
	return &pb.InformeResponse{Recibido: true}, nil	
}

func (s *server) SolicitudPaquete(ctx context.Context, in *pb.RepartoRequest) (*pb.RepartoResponse, error) {
	log.Println("Recibiendo solicitud de paquetes desde el camion", in.IdCamion)
	it, flag := entregarPaq(in.IdCamion)
	if flag{
	// Generamos y retornamos la respuesta de tipo CorreoResponse
		return &pb.RepartoResponse{Id: paquetes[it].Id,
			Tipo:  paquetes[it].Tipo, //0=es retail, 1=pyme, 2=prioritario
			Valor:  paquetes[it].Valor, 
			Origen:  paquetes[it].Origen,
			Destino:  paquetes[it].Destino,
			Intento:  paquetes[it].Intentos}, nil
	}
	return &pb.RepartoResponse{Id: "",
		Tipo:  0, 
		Valor:  0, 
		Origen:  "",
		Destino:  "",
		Intento:  0}, nil
}


func (s *server) EnviarPaquete(ctx context.Context, in *pb.PaqueteRequest) (*pb.PaqueteResponse, error) {
	log.Printf("Recibiendo solicitud de envio paquete con destino a %v",in.Destino)
	n += 1
	cSeg := generarId()
	iPaq := generarId2(n)
	paquetes = append(paquetes, Paquete{
		Id: iPaq, 
		Tipo: in.Tipo, 
		Nombre: in.Nombre, 
		Valor: in.Valor, 
		Origen: in.Origen, 
		Destino: in.Destino, 
		Seguimiento: cSeg,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Estado: "En bodega",
		Intentos: 0,
		Intentar: true})
	
	if in.Tipo == 0{
		colaRetail = append(colaRetail, iPaq)
		log.Println("El paquete recibido es del tipo Retail")
		//log.Println("")
		//log.Println(colaRetail)
	}else if in.Tipo == 1{
		colaNormal = append(colaNormal, iPaq)
		log.Println("El paquete recibido es del tipo Pyme")
		//log.Println("")
		//log.Println(colaNormal)
		
	}else if in.Tipo == 2{
		colaPrioritario = append(colaPrioritario, iPaq)
		log.Println("El paquete recibido es del tipo Prioritario")
		//log.Println("")
		//log.Println(colaPrioritario)
	}
	log.Println("El codigo de seguimiento generado es: ",cSeg)
	// Generamos y retornamos la respuesta de tipo PaqueteResponse
	return &pb.PaqueteResponse{NumSeguimiento: cSeg}, nil
	

}

func (s *server) SeguimientoPaquete (ctx context.Context, in *pb.SeguimientoRequest) (*pb.SeguimientoResponse, error) {
	log.Printf("Solicitud de seguimiento de paquete %v recibida", in.NumSeguimiento)
	var est string
	for i := 0; i < len(paquetes); i++ {
		if paquetes[i].Seguimiento ==  in.NumSeguimiento{
			est = paquetes[i].Estado
		}
	}
	return &pb.SeguimientoResponse{Estado: est}, nil
}

func servicio1(){
	// Exponemos nuestro servidor para que reciba llamados
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPaqueteServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func servicio2(){
	// Exponemos nuestro servidor para que reciba llamados
	lis, err := net.Listen("tcp", port2)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPaqueteServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func checkError(message string, err error) {
    if err != nil {
        log.Fatal(message, err)
    }
}

func printPaqs(){
	start := time.Now()
	for{
		if time.Since(start).String()[:2] == "1m"{
			time.Sleep(10 * time.Second)
		    file, err := os.Create("historico.csv")
		    checkError("Cannot create file", err)
		    defer file.Close()

		    writer := csv.NewWriter(file)
		    defer writer.Flush()

		    for i := 0; i < len(paquetes); i++ {
		    	s := make([]string, 0)
		    	if i == 0 {
		    	//	
		    		err1 := writer.Write([]string{"Id paquete","Tipo","Nombre producto","Valor","Origen","Destino","Numero seguimiento"	,"Fecha ingreso","Estado paquete","Intentos"})
		    		checkError("Cannot write to file", err1)
		    	}
		   		s = append(s, paquetes[i].Id)
		   		if paquetes[i].Tipo == 0{
		   			s = append(s, "Retail")
		   		}else if paquetes[i].Tipo == 1{
		   			s = append(s, "Normal")
		   		}else{
		   			s = append(s, "Prioritario")
		   		}
		    	s = append(s, paquetes[i].Nombre)
		    	a := int(paquetes[i].Valor)
		    	s = append(s, strconv.Itoa(a))
		    	s = append(s, paquetes[i].Origen)
		    	s = append(s, paquetes[i].Destino)
		    	s = append(s, paquetes[i].Seguimiento)
		    	s = append(s, paquetes[i].Timestamp)
		    	s = append(s, paquetes[i].Estado)
		    	a = int(paquetes[i].Intentos)
		    	s = append(s, strconv.Itoa(a))
		    	//s = append(s, strconv.FormatBool(paquetes[i].Intentar))
		    	err := writer.Write(s)
		        checkError("Cannot write to file", err)
		    }
			break
		}
	}
	
	log.Println("Pasaron aproximadamente 1[min] y 10[s] . Se esta finalizando la simulacion...")
	time.Sleep(2*time.Second)
	log.Println("Presionar enter para salir")
	//log.Fatal("_____________________________________________________________________________")
}




func main() {
	go servicio1()
	go servicio2()
	go printPaqs()
	fmt.Scanln()

}