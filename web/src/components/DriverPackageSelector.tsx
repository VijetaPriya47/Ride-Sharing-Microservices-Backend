import { PackagesMeta } from './PackagesMeta'
import { CarPackageSlug } from '../types'
import { cn } from "../lib/utils"
import { Check } from 'lucide-react'

interface DriverPackageSelectorProps {
  onSelect: (packageSlug: CarPackageSlug) => void
}

export function DriverPackageSelector({ onSelect }: DriverPackageSelectorProps) {
  return (
    <div className="flex items-center justify-center min-h-screen bg-slate-50/50 p-4">
      <div className="w-full max-w-4xl mx-auto">
        <div className="text-center mb-10">
          <h2 className="text-3xl font-bold text-slate-900 mb-3 tracking-tight">Select Your Vehicle Class</h2>
          <p className="text-slate-500 text-lg">Choose the service level you want to provide.</p>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {Object.entries(PackagesMeta).map(([slug, meta]) => (
            <div
              key={slug}
              className={cn(
                "group relative flex flex-col p-6 rounded-2xl border-2 border-slate-100 bg-white transition-all duration-300 cursor-pointer",
                "hover:border-slate-300 hover:shadow-xl hover:shadow-slate-200/50 hover:-translate-y-1",
                "active:scale-[0.98]"
              )}
              onClick={() => onSelect(slug as CarPackageSlug)}
            >
              <div className="flex items-start justify-between mb-4">
                <div className={cn(
                  "p-3 rounded-xl bg-slate-50 text-slate-900 transition-colors group-hover:bg-slate-900 group-hover:text-white",
                )}>
                  {meta?.icon}
                </div>
                <div className="opacity-0 group-hover:opacity-100 transition-opacity">
                  <div className="w-6 h-6 rounded-full bg-slate-900 text-white flex items-center justify-center">
                    <Check className="w-3.5 h-3.5" />
                  </div>
                </div>
              </div>

              <div className="mt-auto">
                <h3 className="font-bold text-lg text-slate-900 mb-1">{meta?.name}</h3>
                <p className="text-sm text-slate-500 leading-relaxed">{meta?.description}</p>
              </div>
            </div>
          ))}
        </div>
        
        <div className="mt-10 text-center">
          <button 
            onClick={() => window.location.reload()}
            className="text-slate-400 hover:text-slate-600 text-sm font-medium transition-colors"
          >
            ← Back to selection
          </button>
        </div>
      </div>
    </div>
  )
}
