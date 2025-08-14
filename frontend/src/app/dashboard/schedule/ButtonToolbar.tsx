'use client'
import { Button } from '@/components/ui/Button'
import { ButtonGroup } from '@/components/ui/ButtonGroup'

export default function ButtonToolbar({ onReduce, onAuto, onAI }: { onReduce: () => void; onAuto: () => void; onAI: () => void }) {
  return (
    <ButtonGroup direction="horizontal" spacing="normal">
      <Button onClick={onReduce}>ปรับลดพนักงาน</Button>
      <Button onClick={onAuto}>สร้างอัตโนมัติ (Backend)</Button>
      <Button onClick={onAI}>สร้างอัตโนมัติ (AI)</Button>
    </ButtonGroup>
  )
}


