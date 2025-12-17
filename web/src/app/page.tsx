"use client"

// Assets
import 'leaflet/dist/leaflet.css';
// Fix for default marker icon
import icon from 'leaflet/dist/images/marker-icon.png'
import iconShadow from 'leaflet/dist/images/marker-shadow.png'
import dynamic from 'next/dynamic'
import { Button } from "../components/ui/button";
import { useState, Suspense } from "react";
import { useSearchParams, useRouter } from 'next/navigation';
import { CarPackageSlug } from '../types';
import { DriverPackageSelector } from '../components/DriverPackageSelector';
import { Car } from 'lucide-react';

// Dynamically import components that use Leaflet
const DriverMap = dynamic(() => import("../components/DriverMap").then(mod => mod.DriverMap), { ssr: false })
const RiderMap = dynamic(() => import("../components/RiderMap"), { ssr: false })

// Initialize Leaflet icon only on client side
if (typeof window !== 'undefined') {
  import('leaflet').then((L) => {
    const DefaultIcon = L.default.icon({
      iconUrl: icon.src,
      shadowUrl: iconShadow.src,
      iconSize: [25, 41],
      iconAnchor: [12, 41],
    })
    L.default.Marker.prototype.options.icon = DefaultIcon
  })
}

function HomeContent() {
  const [userType, setUserType] = useState<"driver" | "rider" | null>(null)
  const router = useRouter()
  const searchParams = useSearchParams()
  const payment = searchParams.get("payment")
  const [packageSlug, setPackageSlug] = useState<CarPackageSlug | null>(null)

  const handleClick = (userType: "driver" | "rider") => {
    setUserType(userType)
  }

  if (payment === 'success') {
    return (
      <main className="min-h-screen bg-slate-50 flex items-center justify-center p-4">
        <div className="bg-white/80 backdrop-blur-xl p-8 rounded-3xl shadow-2xl border border-white/20 max-w-md w-full text-center">
          <div className="mb-8">
            <div className="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6 shadow-inner">
              <svg className="w-10 h-10 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2.5" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Payment Successful!</h1>
            <p className="text-slate-500 mt-3 text-lg">Your ride has been confirmed and is on the way.</p>
          </div>
          <Button
            className="w-full text-lg py-7 rounded-xl font-semibold shadow-lg shadow-primary/20 hover:shadow-primary/30 transition-all duration-300"
            variant="default"
            onClick={() => router.push("/")}
          >
            Return Home
          </Button>
        </div>
      </main>
    )
  }

  return (
    <main className="min-h-screen bg-gradient-to-br from-slate-100 via-white to-slate-100">
      {userType === null && (
        <div className="flex flex-col items-center justify-center min-h-screen gap-8 px-4 relative overflow-hidden">
          {/* Decorative background elements */}
          <div className="absolute top-[-20%] right-[-10%] w-[600px] h-[600px] bg-blue-100/50 rounded-full blur-3xl" />
          <div className="absolute bottom-[-20%] left-[-10%] w-[500px] h-[500px] bg-indigo-100/50 rounded-full blur-3xl" />

          <div className="bg-white/70 backdrop-blur-2xl p-8 sm:p-12 rounded-[2rem] shadow-2xl border border-white/50 max-w-lg w-full z-10">
            <div className="text-center mb-10">
              <h2 className="text-4xl font-extrabold text-slate-900 mb-4 tracking-tight">RideShare</h2>
              <p className="text-slate-500 text-lg">Choose how you want to move today.</p>
            </div>
            
            <div className="space-y-4">
              <Button
                className="w-full text-lg h-auto py-6 rounded-xl bg-slate-900 text-white hover:bg-slate-800 shadow-xl shadow-slate-900/10 hover:shadow-slate-900/20 transition-all duration-300 group"
                onClick={() => handleClick("rider")}
              >
                <div className="flex items-center justify-center gap-3">
                  <Car className="w-6 h-6 group-hover:scale-110 transition-transform" />
                  <span className="font-semibold">I Need a Ride</span>
                </div>
              </Button>
              
              <Button
                className="w-full text-lg h-auto py-6 rounded-xl border-2 border-slate-200 bg-white/50 hover:bg-white hover:border-slate-300 text-slate-700 shadow-sm hover:shadow-md transition-all duration-300 group"
                variant="outline"
                onClick={() => handleClick("driver")}
              >
                <div className="flex items-center justify-center gap-3">
                  <Car className="w-6 h-6 group-hover:rotate-12 transition-transform" />
                  <span className="font-semibold">I Want to Drive</span>
                </div>
              </Button>
            </div>
          </div>
        </div>
      )}

      {userType === "driver" && packageSlug && (
        <DriverMap packageSlug={packageSlug} />
      )}

      {userType === "driver" && !packageSlug && (
        <DriverPackageSelector onSelect={setPackageSlug} />
      )}

      {userType === "rider" && <RiderMap />}
    </main>
  );
}

export default function Home() {
  return (
    <Suspense fallback={
      <main className="min-h-screen bg-slate-50 flex items-center justify-center">
        <div className="animate-pulse flex flex-col items-center">
          <div className="h-12 w-12 bg-slate-200 rounded-full mb-4"></div>
          <div className="h-4 w-32 bg-slate-200 rounded"></div>
        </div>
      </main>
    }>
      <HomeContent />
    </Suspense>
  );
}
