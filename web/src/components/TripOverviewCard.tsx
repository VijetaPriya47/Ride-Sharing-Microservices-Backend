import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "./ui/card"
import { ScrollArea } from "./ui/scroll-area"
import { cn } from "../lib/utils"

interface TripOverviewCardProps {
  title: string
  description: string
  children?: React.ReactNode
  className?: string
}

export const TripOverviewCard = ({ title, description, children, className }: TripOverviewCardProps) => {
  return (
    <Card className={cn(
      "w-full md:max-w-[420px] z-[1000] shadow-2xl border-0 ring-1 ring-black/5 bg-white/95 backdrop-blur-md rounded-2xl overflow-hidden flex flex-col max-h-[85vh]", 
      className
    )}>
      <CardHeader className="bg-white/50 border-b border-slate-100/50 pb-4 pt-5 px-6 backdrop-blur-sm sticky top-0 z-10">
        <CardTitle className="text-xl font-bold text-slate-900 tracking-tight">{title}</CardTitle>
        <CardDescription className="text-slate-500 font-medium">{description}</CardDescription>
      </CardHeader>
      <CardContent className="p-0 flex-1 relative overflow-hidden">
        <ScrollArea className="h-full max-h-[calc(85vh-90px)]">
          <div className="p-5">
            {children}
          </div>
        </ScrollArea>
      </CardContent>
    </Card>
  )
}
