import { credentials } from "grpc";
import {LocationServiceClient} from "../../proto/car_location_grpc_pb";
import {CarLocation, CarLocationRequest} from "../../proto/car_location_pb";
import { Empty } from "google-protobuf/google/protobuf/empty_pb";
import { noop } from "../users/utils";

const port = 3000;

export const defaultLocationClient = new LocationServiceClient(
  `localhost:${port}`,
  credentials.createInsecure()
);

export class CarLocationService {
  
   private client: LocationServiceClient
  
   constructor(address: string) {
     this.client = new LocationServiceClient(
       address,
       credentials.createInsecure()
     );
   }
   
   setClient(client: LocationServiceClient){
     this.client = client
   }
   
   
    getAllCarLocation() : Promise<CarLocation[]> {
     return new Promise<CarLocation[]>((resolve, reject) => {
       const stream = this.client.getAllCarLocation(new Empty());
       const locations: CarLocation[] = [];
       stream.on("data", (car) => {
         console.log("received", car)
         locations.push(car)
       });
       stream.on("error", reject);
       stream.on("end", () => resolve(locations));
     });
   }
   
   getAllCarLocationWithCallBack(callback: (car: CarLocation) => void) : Promise<boolean> {
     return new Promise<boolean>((resolve, reject) => {
       const stream = this.client.getAllCarLocation(new Empty());
       stream.on("data", (car) => {
        //  console.log("received", car)
        callback(car)
       });
       stream.on("error", reject);
       stream.on("end", () => resolve(true));
     });
   }
   
   createNewUsers(car: CarLocation, callback: () => void) {
     this.client.submitLocation(car, callback);
   }
   
   getCarLocation(carId: string) : Promise<CarLocation[]> {
     const request = new CarLocationRequest();
     request.setCarid(carId);
     return new Promise<CarLocation[]>((resolve, reject) => {
       const stream = this.client.getCarLocation(request, new Empty());
       const locations: CarLocation[] = [];
       stream.on("data", (car) => locations.push(car));
       stream.on("error", reject);
       stream.on("end", () => resolve(locations));
     });
   }
   
}
