import { Button } from "./ui/button"
import { Clock } from 'lucide-react'
import { RouteFare, TripPreview } from '../types'
import { convertMetersToKilometers, convertSecondsToMinutes } from "../utils/math"
import { cn } from "../lib/utils"
import { PackagesMeta } from "./PackagesMeta"
import { TripOverviewCard } from "./TripOverviewCard"

interface DriverListProps {
  trip: TripPreview | null;
  onPackageSelect: (fare: RouteFare) => void
  onCancel: () => void
}


export function DriverList({ trip, onPackageSelect, onCancel }: DriverListProps) {
  return (
    <TripOverviewCard
      title="Select Ride"
      description={`Trip distance: ${convertMetersToKilometers(trip?.distance ?? 0)}`}
    >
      <div className="flex items-center gap-2 text-sm text-slate-500 mb-6 bg-slate-50 p-3 rounded-lg border border-slate-100">
        <Clock className="w-4 h-4 text-slate-400" />
        <span className="font-medium">Estimated arrival: {convertSecondsToMinutes(trip?.duration ?? 0)}</span>
      </div>
      
      <div className="space-y-3">
        {trip?.rideFares.map((fare) => {
          const Icon = PackagesMeta[fare.packageSlug].icon;
          const price = fare.totalPriceInCents && `$${(fare.totalPriceInCents / 100).toFixed(2)}`

          return (
            <div
              key={fare.id}
              className={cn(
                "group relative flex items-center justify-between p-4 rounded-xl border border-slate-100 bg-white transition-all duration-200 cursor-pointer",
                "hover:border-slate-300 hover:shadow-md hover:translate-x-1",
                "active:scale-[0.99] active:border-slate-400"
              )}
              onClick={() => onPackageSelect(fare)}
            >
              <div className="flex items-center gap-4">
                <div className="p-2.5 bg-slate-50 rounded-lg text-slate-700 group-hover:bg-slate-900 group-hover:text-white transition-colors">
                  {Icon}
                </div>
                <div>
                  <h3 className="font-bold text-slate-900">{PackagesMeta[fare.packageSlug].name}</h3>
                  <p className="text-xs font-medium text-slate-500">{PackagesMeta[fare.packageSlug].description}</p>
                </div>
              </div>
              <div className="text-right">
                <p className="font-bold text-lg text-slate-900">{price}</p>
              </div>
            </div>
          );
        })}
      </div>
      
      <div className="mt-6">
        <Button
          variant="outline"
          className="w-full border-slate-200 hover:bg-slate-50 text-slate-600"
          onClick={() => onCancel()}
        >
          Cancel & Return to Map
        </Button>
      </div>
    </TripOverviewCard>
  )
}
