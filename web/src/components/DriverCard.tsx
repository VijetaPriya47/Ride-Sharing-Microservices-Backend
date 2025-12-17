import { Driver, CarPackageSlug } from "../types";
import Image from "next/image";
import { Card, CardContent } from "./ui/card";
import { Car, Star, ShieldCheck } from "lucide-react";

export const DriverCard = ({ driver, packageSlug }: { driver?: Driver | null, packageSlug?: CarPackageSlug }) => {
  if (!driver) return null;

  return (
    <Card className="border-0 shadow-xl bg-white/90 backdrop-blur-md overflow-hidden ring-1 ring-black/5">
      <CardContent className="p-0">
        <div className="p-4 flex items-center gap-4">
          <div className="relative">
            {driver.profilePicture ? (
              <Image
                className="rounded-full object-cover ring-2 ring-white shadow-md"
                src={driver.profilePicture}
                alt={`${driver.name}'s profile picture`}
                width={64}
                height={64}
              />
            ) : (
              <div className="w-16 h-16 rounded-full bg-slate-100 flex items-center justify-center ring-2 ring-white shadow-md">
                <span className="text-xl font-bold text-slate-400">{driver.name.charAt(0)}</span>
              </div>
            )}
            <div className="absolute -bottom-1 -right-1 bg-white rounded-full p-1 shadow-sm">
               <ShieldCheck className="w-4 h-4 text-green-500 fill-green-50" />
            </div>
          </div>
          
          <div className="flex-1 min-w-0">
            <h3 className="text-lg font-bold text-slate-900 truncate">{driver.name}</h3>
            <div className="flex items-center gap-2 text-slate-500 text-sm mt-0.5">
               <span className="flex items-center gap-1">
                 <Star className="w-3.5 h-3.5 text-yellow-400 fill-yellow-400" />
                 <span className="font-medium text-slate-700">4.9</span>
               </span>
               <span className="w-1 h-1 rounded-full bg-slate-300" />
               <span className="truncate">{driver.carPlate}</span>
            </div>
          </div>
        </div>

        {(packageSlug || driver.carPlate) && (
          <div className="bg-slate-50/80 px-4 py-3 border-t border-slate-100 flex items-center justify-between text-sm">
            {packageSlug && (
               <div className="flex items-center gap-1.5 text-slate-600">
                  <Car className="w-4 h-4" />
                  <span className="font-medium capitalize">{packageSlug}</span>
               </div>
            )}
            <span className="text-xs font-mono bg-white border border-slate-200 px-2 py-0.5 rounded text-slate-500">
               {driver.carPlate.toUpperCase()}
            </span>
          </div>
        )}
      </CardContent>
    </Card>
  )
};
